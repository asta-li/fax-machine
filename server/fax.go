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
func uploadAndFaxFile(file *multipart.File, faxNumber string) error {
    log.Println("Handling fax request")
    // TODO(asta): Improve error handling.

    // Store file in GCS.
    bucketName := os.Getenv("BUCKET_NAME")
    fileName := guuid.New()
    if err := uploadGCS(file, bucketName, fileName.String()); err != nil {
        return err
    }
    fileUrl := "gs://" + bucketName + "/" + fileName.String()
    log.Println("Uploaded file to", fileUrl)

    // Fax the file.
    if err := faxUploadedFile(fileUrl, faxNumber); err != nil {
        return err
    }

    // Delete the file from GCS.
    if err := deleteGCS(bucketName, fileName.String()); err != nil {
        return err
    }
    log.Println("Deleted file at", fileUrl)

    return nil
}

// Uploads an object to GCS.
// See https://cloud.google.com/storage/docs/uploading-objects#storage-upload-object-go.
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

// Fax a file via GCS.
func faxUploadedFile(fileUrl string, sendToNumber string) error {
    // Set Twilio account credentials.
    accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
    authToken := os.Getenv("TWILIO_AUTH_TOKEN")

    // Populate message data.
    msgData := url.Values{}
    msgData.Set("To", sendToNumber)
    msgData.Set("From", "+12184801688")
    msgData.Set("MediaUrl", fileUrl)
    msgDataReader := *strings.NewReader(msgData.Encode())

    // Format and send the fax request.
    client := &http.Client{}
    urlStr := "https://fax.twilio.com/v1/Faxes"
    req, err := http.NewRequest("POST", urlStr, &msgDataReader)
    if err != nil {
        return err
    }
    req.SetBasicAuth(accountSid, authToken)
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Handle the response.
    if (resp.StatusCode < 200 || resp.StatusCode >= 300) {
        fmt.Println("Post request to Twilio unsuccessful", resp.Status);
        // TODO(asta): Return error.
    }

    var data map[string]interface{}
    decoder := json.NewDecoder(resp.Body)
    if err := decoder.Decode(&data); err != nil {
        return err
    }
    fmt.Println("StatusCode", resp.StatusCode)
    fmt.Println("Twilio response status", data["status"])
    fmt.Println("Twilio response message", data["message"])
    fmt.Println("Twilio response", data)
    return nil
}

