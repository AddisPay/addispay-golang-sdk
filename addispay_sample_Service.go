package addispay

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type (
	CheckoutForm struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	}

	PaymentTransaction struct {
		TransactionID string            `json:"transaction_id"`
		User          *User             `json:"user"`
		Amount        float64           `json:"amount"`
		Currency      string            `json:"currency"`
		MerchantFee   float64           `json:"merchant_fee"` // txn fee
		Status        TransactionStatus `json:"status"`
		TxnDate       time.Time         `json:"transaction_date"`
	}

	TransactionList struct {
		Transactions []*PaymentTransaction `json:"transactions"`
		Page         int                    `json:"page"`
		PageSize     int                    `json:"page_size"`
		Total        int                    `json:"total"`
	}

	TransactionStatus string

	User struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
)

const (
	FailedTransactionStatus  TransactionStatus = "failed"
	PendingTransactionStatus TransactionStatus = "pending"
	SuccessTransactionStatus TransactionStatus = "success"
)

// Placeholder data (replace with DB interactions)
var (
	// Mock users
	users = []*User{
		{
			ID:        1002,
			FirstName: "Jon",
			LastName:  "Do",
			Email:     RandomString(5) + "@gmail.com",
		},
		{
			ID:        1032,
			FirstName: "Mary",
			LastName:  "Josef",
			Email:     RandomString(5) + "@gmail.com",
		},
	}

	// Mock transactions
	transactions = []*PaymentTransaction{
		{
			TransactionID: RandomString(10),
			Amount:        10.00,
			MerchantFee:   0.35,
			Currency:      "ETB",
			TxnDate:       time.Now(),
			User:          users[0],
		},
		{
			TransactionID: RandomString(10),
			Amount:        120.00,
			MerchantFee:   1.35,
			Currency:      "USD",
			TxnDate:       time.Now(),
			User:          users[1],
		},
	}
)

type (
	AddispayPaymentService interface {
		Checkout(ctx context.Context, userID int64, form *CheckoutForm) (*PaymentTransaction, error)
		ListPaymentTransactions(ctx context.Context, page, pageSize int) (*TransactionList, error)
	}

	AppAddispayPaymentService struct {
		mu                     *sync.Mutex
		paymentGatewayProvider API
	}
)

// NewAddispayPaymentService constructor
func NewAddispayPaymentService(paymentGatewayProvider API) *AppAddispayPaymentService {
	return &AppAddispayPaymentService{
		mu:                     &sync.Mutex{},
		paymentGatewayProvider: paymentGatewayProvider,
	}
}

// Checkout processes a payment request
func (s *AppAddispayPaymentService) Checkout(ctx context.Context, userID int64, form *CheckoutForm) (*PaymentTransaction, error) {

	// Validate the checkout form
	if err := s.validateCheckoutForm(form); err != nil {
		return nil, err
	}

	// Fetch user by ID
	user, err := s.userByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create a payment request
	invoice := &PaymentRequest{
		Amount:         form.Amount,
		Currency:       form.Currency,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		CallbackURL:    os.Getenv("ADDISPAY_CALLBACK_URL"), // Now environment variable
		TransactionRef: RandomString(10),
	}

	// Send request to Addispay
	response, err := s.paymentGatewayProvider.PaymentRequest(invoice)
	if err != nil {
		return nil, err
	}

	if response.Status != "success" {
		log.Printf("[ERROR] Failed to checkout user request response = [%+v]", response)
		return nil, fmt.Errorf("failed to checkout: %v", response.Message)
	}

	transaction := &PaymentTransaction{
		TransactionID: invoice.TransactionRef,
		Amount:        form.Amount,
		Currency:      form.Currency,
		User:          user,
		Status:        PendingTransactionStatus,
		TxnDate:       time.Now(),
	}

	// Save the transaction
	err = s.savePaymentTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// ListPaymentTransactions lists transactions with pagination
func (s *AppAddispayPaymentService) ListPaymentTransactions(ctx context.Context, page, pageSize int) (*TransactionList, error) {
	// Simple pagination logic
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(transactions) {
		return &TransactionList{
			Transactions: []*PaymentTransaction{},
			Page:         page,
			PageSize:     pageSize,
			Total:        len(transactions),
		}, nil
	}

	if end > len(transactions) {
		end = len(transactions)
	}

	transactionList := &TransactionList{
		Transactions: transactions[start:end],
		Page:         page,
		PageSize:     pageSize,
		Total:        len(transactions),
	}

	return transactionList, nil
}

// savePaymentTransaction saves the transaction in a thread-safe manner
func (s *AppAddispayPaymentService) savePaymentTransaction(ctx context.Context, transaction *PaymentTransaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	transactions = append([]*PaymentTransaction{transaction}, transactions...)

	return nil
}

// userByID fetches user by ID (placeholder)
func (s *AppAddispayPaymentService) userByID(ctx context.Context, userID int64) (*User, error) {
	for _, user := range users {
		if user.ID == userID {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// validateCheckoutForm validates the checkout form
func (s *AppAddispayPaymentService) validateCheckoutForm(form *CheckoutForm) error {
	if form.Amount <= 0 {
		return errors.New("invalid amount: must be greater than 0")
	}
	if form.Currency == "" {
		return errors.New("currency is required")
	}
	// Additional validations like currency checks can be added
	return nil
}
