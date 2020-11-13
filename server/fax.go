package main

import (
    "encoding/json"
	"fmt"
	guuid "github.com/google/uuid"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Handle fax requests.
func uploadAndFaxFile(file *multipart.File, faxNumber string) (string, error) {
	// Store file in GCS.
	fileName := guuid.New()
	if err := uploadGCS(file, fileName.String()); err != nil {
		return "", err
	}

    // Get signed URL for the uploaded object.
    signedUrl, err := getSignedUrl(fileName.String())
	if err != nil {
		return "", err
	}
	log.Println("Uploaded file to", signedUrl)

	// Fax the file.
	faxId, err := faxUploadedFile(signedUrl, faxNumber)
	if err != nil {
		return "", err
	}
    log.Println("Fax ID: ", faxId)
	return faxId, nil
}

// Fax status query response data structure.
type FaxStatusQueryResponse struct {
    Data  FaxStatusResponseData   `json:"data"`
}

type FaxStatusResponseData struct {
    Status        string `json:"status"`
    FailureReason        string `json:"failure_reason,omitempty"`
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
	if (msg.Data.EventType == "fax.delivered" || msg.Data.EventType == "fax.failed") {
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
    Data  FaxSendResponseData   `json:"data"`
}

type FaxSendResponseData struct {
    FaxId  string `json:"id"`
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
