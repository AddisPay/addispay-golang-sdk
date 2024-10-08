# Addispay-Go

This the official Golang SDk for Addispay Payment Gateway

### Usage

##### 1. Installation

```
    go get github.com/AddisPay/addispay-golang-sdk
```

###### API_KEY

Add your `API_KEY: ADDISPAY_xxxxxxxxxx` inside `config.yaml` file.
If you want to run the githb action on your forked repository, you have to create a secrete key named `API_KEY`.

##### 2. Setup

```go
    package main

    import (
        Addispay "github.com/AddisPay/addispay-golang-sdk"
    )

    func main(){
        AddispayAPI := Addispay.NewAddispayAPI()
    }
```

##### 3. Accept Payments

```go
    request := &AddispayAPI.PaymentRequest{
        Amount:         10,
        Currency:       "ETB",
        FirstName:      "Addispay",
        LastName:       "ET",
        Email:          "Addispay@et.io",
        CallbackURL:    "https://posthere.io/e631-44fe-a19e",
        TransactionRef: RandomString(20),
        Customization: map[string]interface{}{
            "title":       "A Unique Title",
            "description": "This a perfect description",
            "logo":        "https://your.logo",
        },
    }

    response, err := AddispayAPI.PaymentRequest(request)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Printf("payment response: %+v\n", response)
```

##### 4. Verify Payment Transactions

```go
    response, err := AddispayAPI.Verify("your-txn-ref")
    if err != nil {
         fmt.Println(err)
         os.Exit(1)
    }

    fmt.Printf("verification response: %+v\n", response)
```

##### 5. Transfer to bank

```go
    request := &BankTransfer{
     AccountName:     "Yinebeb Tariku", 
     AccountNumber:   "34264263", 
     Amount:          10,
     BeneficiaryName: "Yinebeb Tariku",
     Currency:        "ETB",
     Reference:       "3264063st01",
     BankCode:        "32735b19-bb36-4cd7-b226-fb7451cd98f0",
 }
  
 response, err := AddispayAPI.TransferToBank(request)
 fmt.Printf("transfer response: %+v\n", response)
```

### Resources

- <https://developer.Addispay.et/docs/overview/>

### Quirks

Suggestions on how to improve the API:

- Introduction of `status codes` would be a nice to have in the future. Status codes are better than the `message` in a way considering there are so many reasons a transaction could fail.
e.g

```shell
    1001: Success
    4001: DuplicateTransaction
    4002: InvalidCurrency
    4003: InvalidAmount 
    4005: InsufficientBalance  
    5000: InternalServerError
    5001: GatewayError
    5002: RejectedByGateway
```

