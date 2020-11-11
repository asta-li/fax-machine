package main

import (
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/plutov/paypal/v3"
	"log"
	"net/http"
	"os"
	"regexp"
)

var paypalClient *paypal.Client

// Contains file fax metadata.
type FaxResponse struct {
    FaxId string
	Price float32
}

// Upload pdf
func uploadHandlerGin(c *gin.Context) {

	log.Println("Upload pdf")

	// TODO: Perform file validation.
	file, header, _ := c.Request.FormFile("file")
	log.Println("Filename:", header.Filename)
	faxNumber := c.Request.PostFormValue("faxNumber")
	log.Println("Destination fax number:", faxNumber)

	// Upload and fax the file.
	redirectUrl, err := uploadFileAndCreateOrder(&file, faxNumber)
	if err != nil {
		panic(err) // TODO: return proper http error code and fix error handling
	}

	log.Println("Sending upload status response")
	c.JSON(http.StatusOK, gin.H{
		"uploadSuccess": true,
		"redirectUrl": redirectUrl,
	})
}

// Handle fax requests.
func faxHandler(c *gin.Context) {

	log.Println("Handling fax request")

	transactionId := c.Request.PostFormValue("transactionId")

	// perform validation on input
	var re = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)
	if !re.MatchString(transactionId) {
		log.Println("main.faxHandler transaction id format is incorrect, must be an uuid")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "transactionId format is invalid, must be an uuid",
		})
		return
	}

	log.Println("TransactionId:", transactionId)

	// Upload and fax the file.
	faxId, err := faxFile(transactionId)
	if err != nil {
		// TODO: handle this with correct error code
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "something unexpected happened when trying to fax. Please try again later",
		})
		log.Println(err)
		return
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
    FaxId  string `form:"id" json:"id"`
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
    Data  FaxStatusWebhookData   `json:"data"`
}

type FaxStatusWebhookData struct {
    FaxId        string `json:"id"`
    EventType string `json:"event_type"`
    Payload   FaxStatusWebhookPayload `json:"payload"`
}

type FaxStatusWebhookPayload struct {
    FileUrl   string `json:"original_media_url"`
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

	// set up paypal client as a global variable
	client, err := SetUpPaypalClient()
	if err != nil {
		panic(fmt.Errorf("unable to set up paypal client, aborting... %s\n", err))
	}

	paypalClient = client

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
		api.POST("/upload", uploadHandlerGin)
		api.POST("/fax-id", faxHandler)
		api.POST("/fax", faxHandler)
		api.GET("/fax-status", faxQueryHandler)

		api.GET("/paypal", getCreateOrderHandler)
		api.GET("/capture", getCaptureOrderHandler)

		api.GET("/process", getCaptureOrderHandler)
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
