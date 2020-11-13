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
	log.Println("Filename:", header.Filename)
	faxNumber := c.Request.PostFormValue("faxNumber")
	log.Println("Destination fax number:", faxNumber)

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

// TODO(asta): This doesn't work yet.
// Handle fax status webhook.
// Responses: fax.queued, fax.media.processed, fax.sending.started, fax.delivered, fax.failed
// https://developers.telnyx.com/docs/v2/programmable-fax/receiving-webhooks
func faxCompleteHandler(c *gin.Context) {
	log.Println("Handling fax status webhook")
	// TODO(asta): Associate this webhook with a previously sent fax.
	// Note that these hooks may arrive out of order.
	// We'll have to maintain some state about fax IDs and the client will have to poll for results.

	// TODO(asta): How to parse this nested JSON structure? Make JSON data structure.
	responseEvent := c.PostForm("data")
	log.Println("Response:", responseEvent)

	if responseEvent == "fax.delivered" {
		// TODO(asta): The client needs to know that the fax was delivered.
		log.Println("Fax completed!")
	}
}

func main() {

	// GIN
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./client/build", true)))

	// storing up to 5MB in memory.
	router.MaxMultipartMemory = 5 << 20 // 5 MiB
	// Setup route group for the API
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		api.POST("/fax", faxHandlerGin)
	}

	// Fax status and completion webhook.
	router.POST("/fax-complete", faxCompleteHandler)

	router.Run(":3000")

}
