package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/plutov/paypal/v3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//var paypalClient *paypal.Client

var BASE_PRICE float64 = 1.99
var PER_PAGE_PRICE float64 = 1

var orderStruct *paypal.Order //temp for testing

type Credentials struct {
	Paypal struct {
		Account string `account`
		ClientId string	`client_id`
		SecretKey string	`secret_key`
	}
}

// set up the paypal client
func SetUpPaypalClient() (*paypal.Client, error) {

	paypalClientId := os.Getenv("PAYPAL_CLIENT_ID")
	paypalSecretKey := os.Getenv("PAYPAL_SECRET_KEY")

	var credentials Credentials

	// if can't find environment variables for paypal client, then try the yaml file (for dev env)
	if paypalClientId == "" || paypalSecretKey == "" {
		paymentYaml := "./payment.yaml"
		yamlFile, err := ioutil.ReadFile(paymentYaml)
		if err != nil {
			fmt.Printf("Error reading YAML file for payment creds: %s\n", err)
			return nil, err
		}

		err = yaml.Unmarshal(yamlFile, &credentials)
		if err != nil {
			fmt.Printf("Error parsing YAML file: %s\n", err)
			return nil, err
		}

		fmt.Printf("Result: %v\n", credentials.Paypal.Account)

		paypalClientId = credentials.Paypal.ClientId
		paypalSecretKey = credentials.Paypal.SecretKey
	}

	client, err := paypal.NewClient(paypalClientId, paypalSecretKey, paypal.APIBaseSandBox)
	client.SetLog(os.Stdout) // Set log to terminal stdout
	if err != nil {

		return nil, fmt.Errorf("there was an error getting paypal access token at %s, error %s\n", paypal.APIBaseSandBox, err)
	}

	// this is a must do step in order for subsequent requests to paypal to go through
	_, err = client.GetAccessToken()
	if err != nil {
		fmt.Printf("there was an error getting paypal access token at %s, error %s\n", paypal.APIBaseSandBox, err)
		panic(err)
	}

	return client, err
}

func calculateCost(pageCount int) float64 {
	return BASE_PRICE + PER_PAGE_PRICE * float64(pageCount)
}

// sends create order request to Paypal and returns te checkout link to be returned to the front end
func CreatePaypalOrder(pageCount int, metadata TxnMetadata) (string, error) {

	// variables to be passed in
	purchaseValue := fmt.Sprintf("%4.2f", calculateCost(pageCount))
	emailAddress := metadata.EmailAddress

	appContext := paypal.ApplicationContext{
		BrandName:"faxmachine.dev",
		// http://localhost:3000/process?transaction=fbed1651-25a1-4c1c-bb6a-751ac4f613d9&token=10W91560AS205550U&PayerID=ECQGRHBEL4ML6
		CancelURL:fmt.Sprintf("http://localhost:3000/?action=cancel&transaction=%s", metadata.TransactionId),
		ReturnURL:fmt.Sprintf("http://localhost:3000/?action=process&transaction=%s", metadata.TransactionId),
		ShippingPreference:paypal.ShippingPreferenceNoShipping,
	}
	payer := paypal.CreateOrderPayer{EmailAddress:emailAddress}

	var purchaseUnitRequests []paypal.PurchaseUnitRequest
	purchaseUnitRequests = append(purchaseUnitRequests, paypal.PurchaseUnitRequest{
		Description:"fax sent through faxmachine.dev",
		Amount:&paypal.PurchaseUnitAmount{Value: purchaseValue, Currency:"USD"}})

	order, orderErr := paypalClient.CreateOrder(
		paypal.OrderIntentCapture,
		purchaseUnitRequests,
		&payer,
		&appContext,
	)
	if orderErr != nil {
		log.Println(orderErr)
		return "", orderErr
	}
	orderStruct = order

	log.Printf("oder status %s\n", order.Status)
	log.Printf("links from the order %s\n", order.Links)
	log.Printf("links from the order %s\n", order.ID)


	return order.Links[1].Href, nil
}

// after the callback is triggered on the browser, call back to capture the paypal order and also fax the pdf through
func capturePaypalOrder() {

	paypalClient.CaptureOrder(orderStruct.ID, paypal.CaptureOrderRequest{PaymentSource: &paypal.PaymentSource{
		Card:  nil,
		Token: nil,
	}})

}

func getCreateOrderHandler(c *gin.Context) {
	checkoutLink, err := CreatePaypalOrder(4, TxnMetadata{
		TransactionId: "f53920e7-2627-482f-840b-d5fd9d4ac3e2",
		FaxNumber:     "3333333333",
		SignedUrl:     "http://blah",
		EmailAddress:  "sfsdf@sdf.com",
	})

	if err != nil {
		panic(err) // TODO: fix this error handling
	}

	c.JSON(http.StatusOK, gin.H{
		"checkout_link": checkoutLink,
	})
}

func getCaptureOrderHandler(c *gin.Context) {
	capturePaypalOrder()

	c.JSON(http.StatusOK, gin.H{
		"successful": "uncertain",
	})
}

// once user pays through paypal, the paypal will redirect here.
// we will trig
//func getOrderApprovedHandler(c *gin.Context) {
//	file, header, _ := c.Request.FormFile("file")
//	log.Println("Filename:", header.Filename)
//	faxNumber := c.Request.PostFormValue("faxNumber")
//	log.Println("Destination fax number:", faxNumber)

	// Upload and fax the file.
	//redirectUrl, err := uploadFileAndCreateOrder(&file, faxNumber)
//}

func _main() {

	//// set up paypal
	//var err error
	//paypalClient, err = SetUpPaypalClient()
	//if err != nil {
	//	log.Println("FAILED TO CREATE paypal client")
	//}
	//router := gin.Default()
	//api := router.Group("/api")
	//{
	//	api.GET("/paypal", getCreateOrderHandler)
	//	api.GET("/capture", getCaptureOrderHandler)
	//}
	//
	//port := os.Getenv("PORT")
	//if port == "" {
	//	port = "8000"
	//}
	//router.Run(":" + port)

}