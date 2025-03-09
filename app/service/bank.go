package service

import (
	"context"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"sample-bank/app/models"
	pb "sample-bank/proto"
	"time"
)

type BankService struct {
	DB *gorm.DB
	pb.UnimplementedBankServiceServer
}

func (s *BankService) TransferFunds(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if req.FromAccountNumber == req.ToAccountNumber {
		return nil, status.Errorf(codes.InvalidArgument, "cannot transfer same srcAccount")
	}

	srcAccount := models.Accounts{}
	sourceAccount, errSourceAccount := srcAccount.GetUserByAccountNumber(s.DB, req.FromAccountNumber)
	if errSourceAccount != nil {
		return nil, status.Errorf(codes.InvalidArgument, "source srcAccount not found")
	}

	bnfAccount := models.Accounts{}
	beneficiaryAccount, errBeneficiaryAccount := bnfAccount.GetUserByAccountNumber(s.DB, req.ToAccountNumber)
	if errBeneficiaryAccount != nil {
		return nil, status.Errorf(codes.InvalidArgument, "beneficiary srcAccount not found")
	}

	blockedBalance := models.BlockBalances{}
	blockBalance := blockedBalance.GetBlockBalanceByAccountId(s.DB, sourceAccount.Id)

	allowToTransfer := ((sourceAccount.Balance - blockBalance) - req.Amount) >= 0
	if !allowToTransfer {
		return nil, status.Errorf(codes.InvalidArgument, "insuficient balance")
	}

	// create blocked balance
	cb := models.BlockBalances{AccountId: sourceAccount.Id, Amount: req.Amount}
	createBlockBalance, err := cb.CreateBlockBalance(s.DB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	// transaction transfer
	transferId, errTransfer := models.CreateTransaction(s.DB, sourceAccount.Id, sourceAccount.Balance, beneficiaryAccount.Id, beneficiaryAccount.Balance, req.Amount)
	if errTransfer != nil {
		log.Println(errTransfer)
		return nil, status.Errorf(codes.Internal, "Transfers failed id %s", errTransfer.Error())
	}

	// drop blocked balance
	errorDropData := createBlockBalance.DropBlockedBalance(s.DB)
	if errorDropData != nil {
		log.Println(errTransfer)
		return nil, status.Errorf(codes.Internal, "Transfers failed error drop data %s", errorDropData.Error())
	}

	transfer := models.Transfers{}
	transferData, errTf := transfer.GetTransferById(s.DB, transferId)
	if errTf != nil {
		log.Println(errTransfer)
		return nil, status.Errorf(codes.Internal, "Transfers failed data, %s", errTf.Error())
	}

	return &pb.TransferResponse{
		SourceAccountName:        sourceAccount.Name,
		SourceAccountNumber:      sourceAccount.AccountNumber,
		BeneficiaryAccountName:   beneficiaryAccount.Name,
		BeneficiaryAccountNumber: beneficiaryAccount.AccountNumber,
		Amount:                   transferData.Amount,
		TransferStatus:           transfer.Status,
	}, nil
}

func (s *BankService) GetBalance(ctx context.Context, req *pb.BalanceRequest) (*pb.BalanceResponse, error) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	account := models.Accounts{}
	accountData, err := account.GetUserByAccountNumber(s.DB, req.AccountNo)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "account not found")
	}

	return &pb.BalanceResponse{
		AccountName:   accountData.Name,
		AccountNumber: accountData.AccountNumber,
		Balance:       accountData.Balance,
	}, nil
}
