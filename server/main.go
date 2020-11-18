package main

import (
	"fmt"
	"github.com/gin-gonic/contrib/static"
	location "github.com/gin-contrib/location"
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

var defaultPort = "3000"
var defaultHost = "localhost"

// Upload pdf
func uploadHandlerGin(c *gin.Context) {

	log.Println("Upload pdf")
	// TODO: Perform file validation.
	file, header, _ := c.Request.FormFile("file")
	log.Println("Filename:", header.Filename)
	faxNumber := c.Request.PostFormValue("faxNumber")
	log.Println("Destination fax number:", faxNumber)

	// Upload and fax the file.
	redirectUrl, err := uploadFileAndCreateOrder(&file, faxNumber, location.Get(c))
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

	log.Println("main.faxHandler: Handling fax request")

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

	log.Println("main.faxHandler: TransactionId:", transactionId)

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

// Handle fax status queries.
// https://developers.telnyx.com/docs/api/v2/programmable-fax/Programmable-Fax-Commands#ViewFax
func faxQueryHandler(c *gin.Context) {
	log.Println("Handling fax status query")
	faxId := c.DefaultQuery("id", "")

	if faxId == "" {
		log.Println("query parameter fax id : `id` is missing from query parameters")
		c.JSON(http.StatusBadRequest, gin.H{"status": "Unable to parse JSON request"})
		return
	}

	status, err := getFaxStatus(faxId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "something went wrong with getting fax status"})
		return

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

	// print file names and line numbers
	log.SetFlags(log.LstdFlags | log.Lshortfile)



	// set up paypal client as a global variable
	client, err := SetUpPaypalClient()
	if err != nil {
		panic(fmt.Errorf("unable to set up paypal client, aborting... %s\n", err))
	}

	paypalClient = client

	// GIN
	router := gin.Default()
	router.AppEngine = true

	// set up location middleware. Default to localhost & port if it cannot ascertain request host / port
	// https://github.com/gin-contrib/location
	locationConfig := location.DefaultConfig()
	locationConfig.Host = fmt.Sprintf("%s:%s", defaultHost, defaultPort)
	router.Use(location.New(locationConfig))

	// serve the static site with the static middleware
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
		port = defaultPort
	}
	router.Run(":" + port)

}
