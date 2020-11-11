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
  "os"
  "time"
)

// Contains file fax metadata.
type FaxResponse struct {
  Price float32
}


// Uploads an object to GCS.
// See https://cloud.google.com/storage/docs/uploading-objects#storage-upload-object-go.
func storeGCS(dataToWrite multipart.File, bucketName string, fileName string) error {
  // Create GCS connection
  ctx := context.Background()
  client, err := storage.NewClient(ctx)
  if err != nil {
    return fmt.Errorf("storage.NewClient: %v", err)
  }
  defer client.Close()

  // TODO(asta): What is this doing?
  ctx, cancel := context.WithTimeout(ctx, time.Second*50)
  defer cancel()

  // Upload an object with storage.Writer.
  w := client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
  if _, err = io.Copy(w, dataToWrite); err != nil {
    return fmt.Errorf("io.Copy: %v", err)
  }
  if err := w.Close(); err != nil {
    return fmt.Errorf("Writer.Close: %v", err)
  }
  return nil
}

// Handle fax requests.
func faxHandler(w http.ResponseWriter, r *http.Request){
  log.Println("Handling fax request")
  // TODO(asta): Improve error handling.

  // TODO(asta): Perform server-side file validation.
  // Parse form, storing up to 5MB in memory.
  if err := r.ParseMultipartForm(5 << 20); err != nil {
    log.Fatal(err)
  }

  file, header, err := r.FormFile("file")
  defer file.Close()
  log.Println(header.Header)

  // Store file in GCS.
  bucketName := os.Getenv("BUCKET_NAME")
  fileName := guuid.New()
  storeGCS(file, bucketName, fileName.String())
  filePath := "gs://" + bucketName + "/" + fileName.String()
  log.Println("Uploaded file to", filePath)


  // TODO(asta): Fax the file.


  // TODO(asta): Delete the file from GCS.


  // Create successful response data.
  faxResponse := FaxResponse{
    Price: 3.19,
  }

  // Serialize and send the response data.
  faxResponseJson, err := json.Marshal(faxResponse)
  if err != nil{
    panic(err)
  }

  log.Println("Sending fax response")
  w.WriteHeader(http.StatusOK)
  w.Header().Set("Content-Type","application/json")
  w.Write(faxResponseJson)
}

func main() {
    log.Println("Starting server")

    // Serve the static webpage.
    buildHandler := http.FileServer(http.Dir("client/build"))
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        buildHandler.ServeHTTP(w, r)
    })

    // Handle fax requests.
    http.HandleFunc("/fax", faxHandler)

    // Start the server, defaulting to port 3000.
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }
    log.Println("Running server on port", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
}
