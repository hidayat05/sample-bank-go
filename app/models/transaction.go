package models

import (
	"gorm.io/gorm"
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

	updateBalance := sourceAccountBalance - amount
	if err := tx.Table("accounts").
		Where("id = ?", sourceAccountId).
		Update("balance", updateBalance).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	currentBeneficiaryBalance := beneficiaryAccountBalance + amount
	if err := tx.Table("accounts").
		Where("id = ?", beneficiaryId).
		Update("balance", currentBeneficiaryBalance).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	transfer := &Transfers{
		SourceAccountId:      sourceAccountId,
		BeneficiaryAccountId: beneficiaryId,
		Amount:               amount,
		Status:               TransferSuccess,
	}
	if err := tx.Table("transfers").Create(transfer).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return transfer.Id, tx.Commit().Error
}
