package main

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

// Contains file fax metadata.
type FaxResponse struct {
	FaxId string
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

	// Upload and fax the file.
	faxId, err := uploadAndFaxFile(&file, faxNumber)
	if err != nil {
		panic(err)
	}

	// Create successful response data.
	faxResponse := FaxResponse{
		FaxId: faxId,
		Price: 3.19,
	}

	log.Println("Sending fax response")
	c.JSON(http.StatusOK, faxResponse)

}

// Fax status query data structure.
type FaxStatusQuery struct {
	FaxId string `form:"id" json:"id"`
}

// Handle fax status queries.
// https://developers.telnyx.com/docs/api/v2/programmable-fax/Programmable-Fax-Commands#ViewFax
func faxQueryHandler(c *gin.Context) {
	log.Println("Handling fax status query")
	var msg FaxStatusQuery
	if err := c.BindJSON(&msg); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Unable to parse JSON request"})
	}

	status, err := getFaxStatus(msg.FaxId)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

// Fax status webhook data structure.
type FaxStatusWebhook struct {
	Data FaxStatusWebhookData `json:"data"`
}

type FaxStatusWebhookData struct {
	FaxId     string                  `json:"id"`
	EventType string                  `json:"event_type"`
	Payload   FaxStatusWebhookPayload `json:"payload"`
}

type FaxStatusWebhookPayload struct {
	FileUrl string `json:"original_media_url"`
}

// Handle fax status webhook.
// https://developers.telnyx.com/docs/v2/programmable-fax/receiving-webhooks
func faxWebhookHandler(c *gin.Context) {
	log.Println("Handling fax status webhook")
	var msg FaxStatusWebhook
	if err := c.BindJSON(&msg); err != nil {
		log.Fatal(err)
	}

	if err := handleFaxWebhook(msg); err != nil {
		log.Fatal(err)
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
		api.GET("/fax-status", faxQueryHandler)
	}

	// Fax status and completion webhook.
	router.POST("/fax-webhook", faxWebhookHandler)

	// Start the server. Use the environment PORT (e.g. set by Google App Engine),
	// defaulting to port 3000.
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	router.Run(":" + port)

}
