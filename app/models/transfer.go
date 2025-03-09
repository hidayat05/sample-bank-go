package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

type TransferStatus string

const (
	TransferSuccess = "SUCCESS"
	TransferFailed  = "FAILED"
)

type Transfers struct {
	Id                   uint32    `gorm:"primary_key;auto_increment" json:"id"`
	SourceAccountId      uint32    `gorm:"not null" json:"source_account_id"`
	BeneficiaryAccountId uint32    `gorm:"not null" json:"beneficiary_account_id"`
	Amount               float64   `gorm:"not null" json:"amount"`
	Status               string    `gorm:"non null" json:"status"`
	CreatedAt            time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (t *Transfers) CreateTransaction(db *gorm.DB) (*Transfers, error) {
	err := db.Create(&t).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Transfers) GetTransferById(db *gorm.DB, transferId uint32) (*Transfers, error) {
	err := db.Table("transfers").Where("id = ?", transferId).Take(&t).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("transfer data not found")
	}
	return t, nil
}
