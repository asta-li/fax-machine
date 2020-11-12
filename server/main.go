package main

import (
    "cloud.google.com/go/storage"
    "context"
    "fmt"
    "github.com/gin-gonic/contrib/static"
    "github.com/gin-gonic/gin"
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
func faxHandlerGin(c *gin.Context) {

    log.Println("Handling fax request")
    // TODO(asta): Improve error handling.

    // TODO(asta): Perform server-side file validation.

    file, header, _ := c.Request.FormFile("file")
    log.Println(header.Filename)

    //// Upload the file to specific dst.
    //c.SaveUploadedFile(file, dst)


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
    //faxResponseJson, err := json.Marshal(faxResponse)
    //if err != nil{
    //    panic(err)
    //}

    log.Println("Sending fax response")

    c.JSON(http.StatusOK, faxResponse)

    c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", header.Filename))

}

func main() {

    // GIN
    router := gin.Default()
    router.Use(static.Serve("/", static.LocalFile("./client/build", true)))

    // storing up to 5MB in memory.
    router.MaxMultipartMemory = 5 << 20  // 5 MiB
    // Setup route group for the API
    api := router.Group("/api")
    {
        api.GET("/", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H {
                "message": "pong",
            })
        })
        api.POST("/fax", faxHandlerGin)
    }


    router.Run(":3000")


}
