package main

import (
	"encoding/json"
	"fmt"
	guuid "github.com/google/uuid"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// to hold transaction related metadata to be used to continue the transaction after payment
type TxnMetadata struct {
	TransactionId string
	FaxNumber string
	SignedUrl string
	EmailAddress string
}

// Handle fax requests.
func uploadFileAndCreateOrder(file *multipart.File, faxNumber string) (string, error) {
	// Store file in GCS. also used as transactionId

	// TODO: count how many pages the pdf is
	pageCount := 4 // temporary placeholder for pageCount

	fileName := guuid.New()
	log.Println("Attempt to upload file")
	if err := uploadGCS(*file, fileName.String()); err != nil {
		return "", err
	}

	// Get signed URL for the uploaded object.
	signedUrl, err := getSignedUrl(fileName.String())
	if err != nil {
		return "", err
	}

	// write relevant data to a metadata file in GCS to be used to complete the transaction after payment
	txnMetadata := TxnMetadata{
		TransactionId: fileName.String(),
		FaxNumber:     faxNumber,
		SignedUrl:     signedUrl,
		EmailAddress:  "", // TODO: (emma) pass in email?
	}
	err = writeMeta(txnMetadata)
	if err != nil {
		return "", fmt.Errorf("uploadFileAndCreateOrder writeMetadata failure: %v", err)
	}

	log.Println("Uploaded file to", signedUrl)

	//trigger paypal and pass in transaction id.
	redirectUrl, err := CreatePaypalOrder(pageCount, txnMetadata)
	if err != nil {
		log.Printf("error creating paypal order for transaction id %s\n", txnMetadata.TransactionId)
		return "", err
	}

	return redirectUrl, nil
}

// write metadata to GCS used to continue to transaction
// use uuid fileName as the transactionId to be passed back to the front end
func writeMeta(metadata TxnMetadata) error {

	// TODO: emma to delete from disk temp files created
	localMetaFilePath := fmt.Sprintf("./temp/%s.json", metadata.TransactionId)
	gcsMetaFileName := getMetafileName(metadata.TransactionId)

	jsonString, _ := json.Marshal(metadata)

	// if temp directory doesn't exist, create it
	if _, err := os.Stat("./temp"); os.IsNotExist(err) {
		os.Mkdir("./temp", os.ModePerm)
	}

	ioutil.WriteFile(localMetaFilePath,
		jsonString, os.ModePerm)

	// Open local file.
	f, err := os.Open(localMetaFilePath)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	if err := uploadGCS(f, gcsMetaFileName); err != nil {
		log.Println("problem with uploading")
		return err
	}

	log.Println("Uploaded metadata file to", gcsMetaFileName)
	return nil
}

// return the metadata name for given transactionId / fileName (uuid)
func getMetafileName(txnId string) string {
	return fmt.Sprintf("metadata/%s", txnId)

}

func faxFile(transactionId string) (string, error) {

	// TODO: emma need to validate that the person has paid
	// TODO: emma need to check if the person has already faxes this out.
	var metadata TxnMetadata
	// get metadata for transactionId
	data, err := downloadGCS(getMetafileName(transactionId))
	if err != nil {
		return "", fmt.Errorf("download from GCS: %v", err)
	}

	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return "", fmt.Errorf("json unmarshall metadata: %v", err)
	}

	// Fax the file.
	faxId, err := faxUploadedFile(metadata.SignedUrl, metadata.FaxNumber)
	if err != nil {
		return "", err
	}
	log.Println("Fax ID: ", faxId)
	return faxId, nil
}

// Fax status query response data structure.
type FaxStatusQueryResponse struct {
	Data FaxStatusResponseData `json:"data"`
}

type FaxStatusResponseData struct {
	Status        string `json:"status"`
	FailureReason string `json:"failure_reason,omitempty"`
}


// Handle fax status queries.
// https://developers.telnyx.com/docs/api/v2/programmable-fax/Programmable-Fax-Commands#ViewFax
func getFaxStatus(faxId string) (string, error) {
	log.Println("Getting fax status for fax ID:", faxId)

	// Get Telnyx account credentials.
	apiKey := os.Getenv("TELNYX_API_KEY")
	bearer := "Bearer " + apiKey

	// Format and send the fax request.
	client := &http.Client{}
	urlStr := "https://api.telnyx.com/v2/faxes/" + faxId
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	log.Println("Response:", resp)

	// Handle the response.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("Unsuccessful fax status request", resp.Status)
	}

	var data FaxStatusQueryResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&data); err != nil {
		return "", err
	}

	log.Println("Response data:", data)
	if data.Data.Status == "failed" {
		return data.Data.FailureReason, nil
	}
	return data.Data.Status, nil
}

// Handle fax status webhook.
func handleFaxWebhook(msg FaxStatusWebhook) error {
	log.Println("Fax status message:", msg)

	// Delete the media file upon fax completion.
	// Events: fax.queued, fax.media.processed, fax.sending.started, fax.delivered, fax.failed
	if msg.Data.EventType == "fax.delivered" || msg.Data.EventType == "fax.failed" {

		// Signed URLs have the form https://storage.googleapis.com/bucket-name/file-name?signature
		// First, split by '?' to get the non-signed URL
		fileUrl := strings.Split(msg.Data.Payload.FileUrl, "?")[0]
		// First, then split by '/' just to get the object filename.
		fileUrlSplit := strings.Split(fileUrl, "/")
		fileName := fileUrlSplit[len(fileUrlSplit)-1]
		if err := deleteGCS(fileName); err != nil {
			return err
		}
		log.Println("Deleted file at", fileName)
	}
	return nil
}

// Fax status query response data structure.
type FaxSendResponse struct {
	Data FaxSendResponseData `json:"data"`
}

// TODO: figure if there's a way to make initial letter lower case for response data key
type FaxSendResponseData struct {
	FaxId string `json:"id"`
}

// Send a file request to Telnyx Programmable Fax.
func faxUploadedFile(fileUrl string, sendToNumber string) (string, error) {
	// Get Telnyx account credentials.
	sendFromNumber := os.Getenv("FAX_FROM_NUMBER")
	appId := os.Getenv("TELNYX_APP_ID")
	apiKey := os.Getenv("TELNYX_API_KEY")
	bearer := "Bearer " + apiKey

	// Populate message data.
	msgData := url.Values{}
	msgData.Set("connection_id", appId)
	msgData.Set("to", sendToNumber)
	msgData.Set("from", sendFromNumber)
	msgData.Set("media_url", fileUrl)
	msgDataReader := *strings.NewReader(msgData.Encode())

	// Format and send the fax request.
	client := &http.Client{}
	urlStr := "https://api.telnyx.com/v2/faxes"
	req, err := http.NewRequest("POST", urlStr, &msgDataReader)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Handle the response.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("Unsuccessful fax request", resp.Status)
	}

	var data FaxSendResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&data); err != nil {
		return "", err
	}
	log.Println("Response data:", data)
	return data.Data.FaxId, nil
}
