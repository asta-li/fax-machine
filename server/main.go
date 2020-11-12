package main

import (
    "fmt"
    "github.com/gin-gonic/contrib/static"
    "github.com/gin-gonic/gin"
    "log"
    "net/http"
)

// Contains file fax metadata.
type FaxResponse struct {
  Price float32
}

// Handle fax requests.
func faxHandlerGin(c *gin.Context) {

    log.Println("Handling fax request")

    // TODO: Perform file validation.
    file, header, _ := c.Request.FormFile("file")
    log.Println("Filename", header.Filename)
    faxNumber := c.Request.PostFormValue("faxNumber")
    log.Println("Destination fax number", faxNumber)

    //// Upload the file to specific dst.
    //c.SaveUploadedFile(file, dst)

    // Upload and fax the file.
    if err := uploadAndFaxFile(&file, faxNumber); err != nil {
        panic(err)
    }

    // Create successful response data.
    faxResponse := FaxResponse{
        Price: 3.19,
    }

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
