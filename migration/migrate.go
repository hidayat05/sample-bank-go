package migration

import (
	"github.com/jinzhu/gorm"
	"sample-bank/app/models"
)

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(models.Accounts{}, models.Transfers{}, models.BlockBalances{})
	return db
}
