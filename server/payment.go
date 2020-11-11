package main

import (
	"fmt"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/plutov/paypal/v3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

//var paypalClient *paypal.Client

var BASE_PRICE float64 = 1.99
var PER_PAGE_PRICE float64 = 1

var orderStruct *paypal.Order //temp for testing

type Paypal struct {
	Account string `yaml:"PAYPAL_ACCOUNT"`
	ClientId string	`yaml:"PAYPAL_CLIENT_ID"`
	SecretKey string	`yaml:"PAYPAL_SECRET_KEY"`
}

type Credentials struct {
	Paypal Paypal `yaml:"env_variables"`
}

// set up the paypal client
func SetUpPaypalClient() (*paypal.Client, error) {

	paypalClientId := os.Getenv("PAYPAL_CLIENT_ID")
	paypalSecretKey := os.Getenv("PAYPAL_SECRET_KEY")

	var credentials Credentials

	// if can't find environment variables for paypal client, then try the yaml file (for dev env)
	if paypalClientId == "" || paypalSecretKey == "" {

		log.Println("no env variables found for payment credentials, using yaml")

		var paymentYaml string
		isProd := os.Getenv("IS_APPENGINE") != ""
		if isProd {
			paymentYaml =  "./payment.yaml"
			log.Println("using prod credentials for payment")
		} else {
			log.Println("using dev credentials for payment")
			paymentYaml =  "./payment_dev.yaml"
		}

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

		paypalClientId = credentials.Paypal.ClientId
		paypalSecretKey = credentials.Paypal.SecretKey
	}

	fmt.Printf("payment.SetUpPaypalClient: account id for paypal is: %v\n", credentials.Paypal.Account)

	client, err := paypal.NewClient(paypalClientId, paypalSecretKey, paypal.APIBaseSandBox)
	if err != nil {
		return nil, fmt.Errorf("there was an error to setting up the new clientat %s, error %s\n", paypal.APIBaseSandBox, err)
	}
	client.SetLog(os.Stdout) // Set log to terminal stdout

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
func CreatePaypalOrder(pageCount int, metadata TxnMetadata, requestUrl *url.URL) (string, error) {

	// variables to be passed in
	purchaseValue := fmt.Sprintf("%4.2f", calculateCost(pageCount))
	emailAddress := metadata.EmailAddress

//https://0.0.10.44/:5500:c02d:fd00:cc44:26cd:cf2:b6e4,%20169.254.1.1
	// possible valid request origin hosts
	//appengineHost := "fax-machine-295219.wl.r.appspot.com"
	//appengineCustomDomainHost := "2604:5500:c02d:fd00:cc44:26cd:cf2:b6e4, 169.254.1.1"
	//domainCustomHost := "136.24.14.76, 169.254.1.1"
	domainHost := "faxmachine.dev"
	localhost := "localhost"

	var currentHost = requestUrl.Hostname()
	log.Printf("the current request host is: %s\n", currentHost)
	// TODO: figure out a way given appengine's dynamic hostnames to prevent bad actions
	//// if incoming request comes from outside of these hosts, then it's not a valid request. we don't allow
	//// direct access to the API
	//if currentHost != appengineHost && currentHost !=domainHost &&
	//	currentHost != localhost && currentHost != appengineCustomDomainHost &&
	//	currentHost != domainCustomHost {
	//	return "", fmt.Errorf("unknow host: \"%s\" with error \"%s\"",
	//		currentHost, errors.New("request from an unknown host. aborting request"))
	//}

	// if using appengineCustomDomainHost then the traffic is coming from our prod domain
	if currentHost != localhost {
		currentHost = domainHost
	} else {
		currentHost = fmt.Sprintf("%s:%s", defaultHost, defaultPort)
	}

	requestUrl.Host = currentHost
	// reset the query params in case there were any that were passed in (there shouldn't)
	requestUrl.RawQuery = ""

	// construct redirect url from paypal, it usually looks something like
	// http://localhost:3000/process?transaction=fbed1651-25a1-4c1c-bb6a-751ac4f613d9&token=10W91560AS205550U&PayerID=ECQGRHBEL4ML6
	q := requestUrl.Query()
	q.Set("action", "process")
	q.Add("transactionId", metadata.TransactionId)
	requestUrl.RawQuery = q.Encode()
	processUrl := requestUrl.String()

	q.Set("action", "cancel")
	requestUrl.RawQuery = q.Encode()
	cancelUrl := requestUrl.String()

	log.Printf("redirect URL constructed: \n\tprocess URL: %s\n\tcancelURL: %s\n", processUrl, cancelUrl)

	appContext := paypal.ApplicationContext{
		BrandName:"faxmachine.dev",

		CancelURL:cancelUrl,
		ReturnURL:processUrl,
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
	}, location.Get(c))

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