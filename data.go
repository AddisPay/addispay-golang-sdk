package addispay

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	PaymentRequest struct {
		Amount         float64                `json:"amount"`
		Currency       string                 `json:"currency"`
		Email          string                 `json:"email"`
		FirstName      string                 `json:"first_name"`
		LastName       string                 `json:"last_name"`
		Phone          string                 `json:"phone"`
		CallbackURL    string                 `json:"callback_url"`
		TransactionRef string                 `json:"tx_ref"`
		Customization  map[string]interface{} `json:"customization"`
	}

	PaymentResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
		Data    struct {
			CheckoutURL string `json:"checkout_url"`
		}
	}

	VerifyResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
		Data    struct {
			TransactionFee float64 `json:"charge"`
		}
	}
)

type (
	BankTransfer struct {
		AccountName     string  `json:"account_name"`
		AccountNumber   string  `json:"account_number"`
		Amount          float64 `json:"amount"`
		BeneficiaryName string  `json:"beneficiary_name"`
		Currency        string  `json:"currency"`
		Reference       string  `json:"reference"`
		BankCode        string  `json:"bank_code"`
	}

	BankTransferResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
		Data    string `json:"data"`
	}
)

func (p PaymentRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.TransactionRef, validation.Required.Error("transaction reference is required")),
		validation.Field(&p.Currency, validation.Required.Error("currency is required")),
		validation.Field(&p.Amount, validation.Required.Error("amount is required")),
	)
}

func (t BankTransfer) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.AccountName, validation.Required.Error("account name is required")),
		validation.Field(&t.AccountNumber, validation.Required.Error("account number is required")),
		validation.Field(&t.Amount, validation.Required.Error("amount is required")),
		validation.Field(&t.Currency, validation.Required.Error("currency is required")),
		validation.Field(&t.Reference, validation.Required.Error("reference is required")),
		validation.Field(&t.BankCode, validation.Required.Error("bank code is required")),
	)
}
