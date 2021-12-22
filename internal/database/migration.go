package database

import (
	"github.com/jinzhu/gorm"
	"github.com/thefueley/workout-api/internal/comment"
)

// MigrateDb : migrate db and create comment table
func MigrateDB(db *gorm.DB) error {
	if result := db.AutoMigrate(&comment.Comment{}); result.Error != nil {
		return result.Error
	}

	return nil
}
