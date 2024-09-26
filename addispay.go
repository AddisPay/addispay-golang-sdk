package addispay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	acceptPaymentV1APIURL  = "https://checkouts.addispay.et/api/v1/payments/init" // Replace with your API
	verifyPaymentV1APIURL  = "https://checkouts.addispay.et/api/v1/transaction/verify/%v"  // we need to add an api to verify the transaction ref here
	transferToBankV1APIURL = "https://api.addispay.et/api/v1/transfers"              // Replace with your API
)

type API interface {
	PaymentRequest(request *PaymentRequest) (*PaymentResponse, error)
	Verify(txnRef string) (*VerifyResponse, error)
	TransferToBank(request *BankTransfer) (*BankTransferResponse, error)
}

type addispayClient struct {
	apiKey string
	client *http.Client
}

func NewAddispayAPI() API {
	return &addispayClient{
		apiKey: viper.GetString("API_KEY"),
		client: &http.Client{
			Timeout: 1 * time.Minute,
		},
	}
}

func (c *addispayClient) PaymentRequest(request *PaymentRequest) (*PaymentResponse, error) {
	var err error
	if err = request.Validate(); err != nil {
		err := fmt.Errorf("invalid input %v", err)
		log.Printf("error %v input %v", err.Error(), request)
		return &PaymentResponse{}, err
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, acceptPaymentV1APIURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Close = true

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var paymentResponse PaymentResponse

	err = json.Unmarshal(body, &paymentResponse)
	if err != nil {
		return nil, err
	}
	return &paymentResponse, nil
}

func (c *addispayClient) Verify(txnRef string) (*VerifyResponse, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(verifyPaymentV1APIURL, txnRef), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Close = true

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var verifyResponse VerifyResponse

	err = json.Unmarshal(body, &verifyResponse)
	if err != nil {
		return nil, err
	}

	return &verifyResponse, nil
}

func (c *addispayClient) TransferToBank(request *BankTransfer) (*BankTransferResponse, error) {
	var err error
	if err = request.Validate(); err != nil {
		err := fmt.Errorf("invalid input %v", err)
		log.Printf("error %v input %v", err, request)
		return nil, err
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, transferToBankV1APIURL, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("error %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Close = true

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("error %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while reading response body %v", err)
		return nil, err
	}

	response := BankTransferResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error while unmarshaling response %v", err)
		return nil, err
	}
	return &response, nil
}
