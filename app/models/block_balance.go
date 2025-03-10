package models

import "gorm.io/gorm"

type BlockBalances struct {
	Id        uint32  `gorm:"primary_key;auto_increment" json:"id"`
	AccountId uint32  `gorm:"not null" json:"account_id"`
	Amount    float64 `gorm:"not null" json:"amount"`
}

func (b *BlockBalances) TableName() string {
	return "block_balances"
}

type accountBalance struct {
	Amount float64
}

func (b *BlockBalances) CreateBlockBalance(db *gorm.DB) (*BlockBalances, error) {
	err := db.Create(&b).Error
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *BlockBalances) GetBlockBalanceByAccountId(db *gorm.DB, accountId uint32) float64 {
	var balance accountBalance
	err := db.Table("block_balances").Select("sum(amount) as amount").Where("account_id = ?", accountId).Scan(&balance).Error
	if err != nil {
		return 0
	}
	return balance.Amount
}

func (b *BlockBalances) DropBlockedBalance(db *gorm.DB) error {
	return db.Delete(b).Error
}
