package main

import (
	"cloud.google.com/go/storage"
	"context"
    "encoding/json"
	"fmt"
	guuid "github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Handle fax requests.
func uploadAndFaxFile(file *multipart.File, faxNumber string) (string, error) {
	// Store file in GCS.
	bucketName := os.Getenv("BUCKET_NAME")
	fileName := guuid.New()
	if err := uploadGCS(file, bucketName, fileName.String()); err != nil {
		return "", err
	}
	fileUrl := "gs://" + bucketName + "/" + fileName.String()
	log.Println("Uploaded file to", fileUrl)

	// Fax the file.
	faxId, err := faxUploadedFile(fileUrl, faxNumber)
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
    return data.Data.Status, nil
}

// Handle fax status webhook.
func handleFaxWebhook(msg FaxStatusWebhook) error {
    log.Println("Fax status message:", msg)

    // Delete the media file upon fax completion.
    // Events: fax.queued, fax.media.processed, fax.sending.started, fax.delivered, fax.failed
	if (msg.Data.EventType == "fax.delivered" || msg.Data.EventType == "fax.failed") {
	    bucketName := os.Getenv("BUCKET_NAME")
        fileUrlSplit := strings.Split(msg.Data.Payload.FileUrl, "/")
        fileName := fileUrlSplit[len(fileUrlSplit)-1]
	    if err := deleteGCS(bucketName, fileName); err != nil {
            return err
        }
	    log.Println("Deleted file at", fileName)
	}
    return nil
}

// Helper functions
// =============================================================

// Uploads a file to GCS.
// See https://cloud.google.com/storage/docs/uploading-objects#storage-upload-object-go.
// TODO(asta): Update permissions: https://cloud.google.com/storage/docs/access-control/signed-urls
func uploadGCS(dataToWrite *multipart.File, bucketName string, fileName string) error {
	// Create GCS connection
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	w := client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	if _, err = io.Copy(w, *dataToWrite); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

// Write a (key, value) pair associated with the given faxId to GCS.
func writeState(bucketName string, faxId string, key string, value string) error {
	// TODO(asta): Pass in the client and context instead of regenerating each time.
	// Create GCS connection
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	fileName := faxId + "/" + key
	w := client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	if _, err := fmt.Fprintf(w, value); err != nil {
		return fmt.Errorf("fmt.Fprintf to object: %v", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

// Deletes an object from GCS.
func deleteGCS(bucketName string, fileName string) error {
	// Create GCS connection
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	obj := client.Bucket(bucketName).Object(fileName)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", fileName, err)
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
	// TODO(asta): Use the user-given file URL after handling permissions.
	// msgData.Set("MediaUrl", fileUrl)
	msgData.Set("media_url", "https://www.twilio.com/docs/documents/25/justthefaxmaam.pdf")
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

    log.Println("Response:", resp)

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
