package models

import (
	"gorm.io/gorm"
	"time"
)

type Accounts struct {
	Id            uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name          string    `gorm:"size:255;not null" json:"name"`
	AccountNumber string    `gorm:"size:255;unique, not null" json:"account_number"`
	Balance       float64   `gorm:"not null" json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (a *Accounts) TableName() string {
	return "accounts"
}

func (a *Accounts) GetUserByAccountNumber(db *gorm.DB, accountNumber string) (*Accounts, error) {
	err := db.Where("account_number = ?", accountNumber).Take(&a).Error
	if err != nil {
		return &Accounts{}, err
	}
	return a, nil
}

func (a *Accounts) GetUserBalance(db *gorm.DB, accountNumber string) (float64, error) {
	account, err := a.GetUserByAccountNumber(db, accountNumber)
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}

func (a *Accounts) UpdateBalanceByUserId(db *gorm.DB, balance float64) (*Accounts, error) {
	err := db.Model(a).Updates(map[string]interface{}{"balance": balance}).Error
	if err != nil {
		return nil, err
	}
	account, errGetAccount := a.GetUserByAccountNumber(db, a.AccountNumber)
	if errGetAccount != nil {
		return nil, errGetAccount
	}
	return account, nil
}
