package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"

	"sample-bank/app/models"
	pb "sample-bank/proto"
)

// MockAccountModel is a mock for the Accounts model.
type MockAccountModel struct {
	mock.Mock
}

func (m *MockAccountModel) GetUserByAccountNumber(accountNo string) (models.Accounts, error) {
	args := m.Called(accountNo)
	return args.Get(0).(models.Accounts), args.Error(1)
}

func TestTransferFundsSuccess(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.Accounts{}, &models.BlockBalances{}, &models.Transfers{})

	service := &BankService{DB: db}
	mockAccount := new(MockAccountModel)

	srcAccount := models.Accounts{Id: 1, AccountNumber: "123", Balance: 1000, Name: "Aji"}
	bnfAccount := models.Accounts{Id: 2, AccountNumber: "456", Balance: 500, Name: "Ferry"}

	db.Create(&srcAccount)
	db.Create(&bnfAccount)

	mockAccount.On("GetUserByAccountNumber", "123").Return(srcAccount, nil)
	mockAccount.On("GetUserByAccountNumber", "456").Return(bnfAccount, nil)

	req := &pb.TransferRequest{FromAccountNumber: "123", ToAccountNumber: "456", Amount: 100}
	resp, err := service.TransferFunds(context.Background(), req)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Aji", resp.SourceAccountName)
	assert.Equal(t, "Ferry", resp.BeneficiaryAccountName)
}

func TestTransferFundsInsufficientBalance(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.Accounts{}, &models.BlockBalances{}, &models.Transfers{})

	service := &BankService{DB: db}
	mockAccount := new(MockAccountModel)

	srcAccount := models.Accounts{Id: 1, AccountNumber: "123", Balance: 1000, Name: "Aji"}
	bnfAccount := models.Accounts{Id: 2, AccountNumber: "456", Balance: 500, Name: "Ferry"}

	db.Create(&srcAccount)
	db.Create(&bnfAccount)

	mockAccount.On("GetUserByAccountNumber", "123").Return(srcAccount, nil)
	mockAccount.On("GetUserByAccountNumber", "456").Return(bnfAccount, nil)

	req := &pb.TransferRequest{FromAccountNumber: "123", ToAccountNumber: "456", Amount: 11000}
	_, err := service.TransferFunds(context.Background(), req)

	assert.NotNil(t, err)
}

func TestTransferFundsSameAccount(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.Accounts{}, &models.BlockBalances{}, &models.Transfers{})

	service := &BankService{DB: db}
	mockAccount := new(MockAccountModel)

	srcAccount := models.Accounts{Id: 1, AccountNumber: "123", Balance: 1000, Name: "Aji"}

	db.Create(&srcAccount)
	mockAccount.On("GetUserByAccountNumber", "123").Return(srcAccount, nil)

	req := &pb.TransferRequest{FromAccountNumber: "123", ToAccountNumber: "123", Amount: 100}
	_, err := service.TransferFunds(context.Background(), req)
	assert.NotNil(t, err)
}

func TestTransferFundsSourceAccountNotFound(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	_ = db.AutoMigrate(&models.Accounts{}, &models.BlockBalances{}, &models.Transfers{})

	service := &BankService{DB: db}
	mockAccount := new(MockAccountModel)
	srcAccount := models.Accounts{}

	db.Create(&srcAccount)
	mockAccount.On("GetUserByAccountNumber", "123").Return(srcAccount, nil)

	req := &pb.TransferRequest{FromAccountNumber: "123", ToAccountNumber: "245", Amount: 100}
	_, err := service.TransferFunds(context.Background(), req)
	assert.NotNil(t, err)
}

func TestTransferFundsBeneficiaryAccountNotFound(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	_ = db.AutoMigrate(&models.Accounts{}, &models.BlockBalances{}, &models.Transfers{})

	service := &BankService{DB: db}
	mockAccount := new(MockAccountModel)
	srcAccount := models.Accounts{Id: 1, AccountNumber: "123", Balance: 1000, Name: "Aji"}

	db.Create(&srcAccount)
	mockAccount.On("GetUserByAccountNumber", "123").Return(srcAccount, nil)

	req := &pb.TransferRequest{FromAccountNumber: "123", ToAccountNumber: "245", Amount: 100}
	_, err := service.TransferFunds(context.Background(), req)
	assert.NotNil(t, err)
}

func TestGetBalanceSuccess(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.Accounts{})

	service := &BankService{DB: db}

	expectedAccount := &models.Accounts{Name: "Aji", AccountNumber: "123", Balance: 1000}
	db.Create(&expectedAccount)

	mockAccount := &MockAccountModel{}
	mockAccount.On("GetUserByAccountNumber", "123").Return(expectedAccount, nil)

	req := &pb.BalanceRequest{AccountNo: "123"}
	resp, err := service.GetBalance(context.Background(), req)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedAccount.Name, resp.AccountName)
	assert.Equal(t, expectedAccount.Balance, resp.Balance)
}

func TestGetBalance(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	_ = db.AutoMigrate(&models.Accounts{})

	service := &BankService{DB: db}

	expectedAccount := &models.Accounts{}
	db.Create(&expectedAccount)

	mockAccount := &MockAccountModel{}
	mockAccount.On("GetUserByAccountNumber", "123").Return(expectedAccount, nil)

	req := &pb.BalanceRequest{AccountNo: "123"}
	_, err := service.GetBalance(context.Background(), req)

	assert.NotNil(t, err)
}
