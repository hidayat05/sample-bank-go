package models

import (
	"github.com/jinzhu/gorm"
)

func CreateTransaction(
	db *gorm.DB,
	sourceAccountId uint32,
	sourceAccountBalance float64,
	beneficiaryId uint32,
	beneficiaryAccountBalance float64,
	amount float64,
) (uint32, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return 0, err
	}

	sourceAccount := &Accounts{}
	updateBalance := sourceAccountBalance - amount
	if err := tx.Table("accounts").
		Update(sourceAccount).
		Where("id = ?", sourceAccountId).
		Update(map[string]interface{}{"balance": updateBalance}).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	beneficiaryAccount := &Accounts{}
	currentBeneficiaryBalance := beneficiaryAccountBalance + amount
	if err := tx.Table("accounts").
		Update(beneficiaryAccount).
		Where("id = ?", beneficiaryId).
		Update(map[string]interface{}{"balance": currentBeneficiaryBalance}).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	transfer := &Transfers{}
	if err := tx.Table("transfers").Create(&Transfers{
		SourceAccountId:      sourceAccountId,
		BeneficiaryAccountId: beneficiaryId,
		Amount:               amount,
		Status:               TransferSuccess,
	}).Take(transfer).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return transfer.Id, tx.Commit().Error
}
