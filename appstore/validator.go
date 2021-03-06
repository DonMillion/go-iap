package appstore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// SandboxURL is the endpoint for sandbox environment.
	SandboxURL string = "https://sandbox.itunes.apple.com/verifyReceipt"
	// ProductionURL is the endpoint for production environment.
	ProductionURL string = "https://buy.itunes.apple.com/verifyReceipt"
	// ContentType is the request content-type for apple store.
	ContentType string = "application/json; charset=utf-8"
)

// IAPClient is an interface to call validation API in App Store
type IAPClient interface {
	Verify(ctx context.Context, reqBody IAPRequest, resp interface{}) error
}

// Client implements IAPClient
type Client struct {
	ProductionURL string
	SandboxURL    string
	httpCli       *http.Client
	IsProduct     bool
}

// HandleError returns error message by status code
func HandleError(status int) error {
	var message string

	switch status {
	case 0:
		return nil

	case 21000:
		message = "The App Store could not read the JSON object you provided."

	case 21002:
		message = "The data in the receipt-data property was malformed or missing."

	case 21003:
		message = "The receipt could not be authenticated."

	case 21004:
		message = "The shared secret you provided does not match the shared secret on file for your account."

	case 21005:
		message = "The receipt server is not currently available."

	case 21007:
		message = "This receipt is from the test environment, but it was sent to the production environment for verification. Send it to the test environment instead."

	case 21008:
		message = "This receipt is from the production environment, but it was sent to the test environment for verification. Send it to the production environment instead."

	case 21010:
		message = "This receipt could not be authorized. Treat this the same as if a purchase was never made."

	default:
		if status >= 21100 && status <= 21199 {
			message = "Internal data access error."
		} else {
			message = "An unknown error occurred"
		}
	}

	return errors.New(message)
}

// New creates a client object
func New(isProduct bool) *Client {
	client := &Client{
		ProductionURL: ProductionURL,
		SandboxURL:    SandboxURL,
		IsProduct:     isProduct,
		httpCli: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	return client
}

// NewWithClient creates a client with a custom http client.
func NewWithClient(client *http.Client, isProduct bool) *Client {
	return &Client{
		ProductionURL: ProductionURL,
		SandboxURL:    SandboxURL,
		httpCli:       client,
		IsProduct:     isProduct,
	}
}

// Verify sends receipts and gets validation result
func (c *Client) Verify(ctx context.Context, reqBody IAPRequest, result interface{}) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqBody)

	req, err := http.NewRequest("POST", c.ProductionURL, b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", ContentType)
	req = req.WithContext(ctx)
	resp, err := c.httpCli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return c.parseResponse(resp, result, ctx, reqBody)
}

func (c *Client) parseResponse(resp *http.Response, result interface{}, ctx context.Context, reqBody IAPRequest) error {
	// Read the body now so that we can unmarshal it twice
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &result)
	if err != nil {
		return err
	}

	// https://developer.apple.com/library/content/technotes/tn2413/_index.html#//apple_ref/doc/uid/DTS40016228-CH1-RECEIPTURL
	var r StatusResponse
	err = json.Unmarshal(buf, &r)
	if err != nil {
		return err
	}
	if c.IsProduct == false && r.Status == 21007 {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(reqBody)

		req, err := http.NewRequest("POST", c.SandboxURL, b)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", ContentType)
		req = req.WithContext(ctx)
		resp, err := c.httpCli.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}
