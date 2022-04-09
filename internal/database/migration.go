package database

import (
	"github.com/Cingozr/go-rest-api/internal/comment"
	"github.com/jinzhu/gorm"
)

func MigrateDb(db *gorm.DB) error {
	if result := db.AutoMigrate(&comment.Comment{}); result.Error != nil {
		return result.Error
	}
	return nil
}
