package transaction_manager

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/tebrizetayi/ledger_service/internal/storage"
	utils "github.com/tebrizetayi/ledger_service/internal/test_utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddTransaction_NotValidAmount(t *testing.T) {

	// Assign
	testEnv, err := utils.CreateTestEnv()
	if err != nil {
		t.Fatalf("failed to create test env: %v", err)
	}
	defer testEnv.Cleanup()

	storageClient := storage.NewStorageClient(testEnv.DB)
	transactionManager := NewTransactionManagerClient(storageClient)

	user := storage.User{
		ID:      uuid.New(),
		Balance: decimal.NewFromFloat(0),
	}
	err = transactionManager.storageClient.UserRepository.Add(testEnv.Context, user)
	if err != nil {
		t.Fatalf("failed to add user: %v", err)
	}

	// Act
	_, err = transactionManager.AddTransaction(testEnv.Context, Transaction{
		ID:             uuid.New(),
		Amount:         decimal.NewFromFloat(0),
		UserID:         user.ID,
		CreatedAt:      time.Now(),
		IdempotencyKey: uuid.New(),
	})

	// Assert
	assert.Equal(t, ErrInvalidTransaction, err)
}

func TestAddTransaction_Success(t *testing.T) {
	// Assign
	testEnv, err := utils.CreateTestEnv()
	if err != nil {
		t.Fatalf("failed to create test env: %v", err)
	}
	defer testEnv.Cleanup()

	storageClient := storage.NewStorageClient(testEnv.DB)
	transactionManager := NewTransactionManagerClient(storageClient)

	user := storage.User{
		ID:      uuid.New(),
		Balance: decimal.NewFromFloat(0),
	}
	err = transactionManager.storageClient.UserRepository.Add(testEnv.Context, user)
	if err != nil {
		t.Fatalf("failed to add user: %v", err)
	}

	// Act
	transaction, err := transactionManager.AddTransaction(testEnv.Context, Transaction{
		ID:             uuid.New(),
		Amount:         decimal.NewFromFloat(100),
		UserID:         user.ID,
		CreatedAt:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		IdempotencyKey: uuid.New(),
	})
	if err != nil {
		t.Fatalf("failed to add transaction: %v", err)
	}

	// Assert
	assert.True(t, transaction.Amount.Equal(decimal.NewFromFloat(100)), "amount should be 100.0")
	assert.Equal(t, user.ID, transaction.UserID)
}
