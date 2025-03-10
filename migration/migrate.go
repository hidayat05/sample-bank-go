package migration

import (
	"gorm.io/gorm"
	"log"
	"sample-bank/app/models"
)

func DBMigrate(db *gorm.DB) *gorm.DB {
	err := db.AutoMigrate(models.Accounts{}, models.Transfers{}, models.BlockBalances{})
	if err != nil {
		log.Fatalf("failed to migrate accounts: %v", err)
	}
	return db
}
