package main

import (
  "cloud.google.com/go/storage"
  "context"
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "os"
)

// Contains file fax metadata.
type FaxResponse struct {
  Success bool
  Price float32
}

func storeGCS(dataToWrite, bucketName, fileName string) error {
  // Create GCS connection
  ctx := context.Background()
  client, err := storage.NewClient(ctx)
  if err != nil {
    log.Fatal(err)
  }

  // Connect to bucket
  bucket := client.Bucket(bucketName)
  obj := bucket.Object(fileName)

  // Write fileData to obj.
  w := obj.NewWriter(ctx)
  if _, err := fmt.Fprintf(w, dataToWrite); err != nil {
    log.Fatal(err)
  }

  // Close, just like writing a file.
  if err := w.Close(); err != nil {
    log.Fatal(err)
  }
  return nil
}

func faxHandler(w http.ResponseWriter, r *http.Request){
  log.Println("Handling fax request")
  r.ParseForm();

  // TODO(asta): Debug form parsing.
  for key, value := range r.Form {
    fmt.Printf("%s = %s\n", key, value)
    log.Println(key)
    log.Println(value)
  }

  faxResponse := FaxResponse{
    Success: true,
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
