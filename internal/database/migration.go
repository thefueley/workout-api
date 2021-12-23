package database

import (
	"github.com/jinzhu/gorm"
	"github.com/thefueley/workout-api/internal/workout"
)

// MigrateDb : migrate db and create comment table
func MigrateDB(db *gorm.DB) error {
	if workout_result := db.AutoMigrate(&workout.Workout{}); workout_result.Error != nil {
		return workout_result.Error
	}

	if exercise_result := db.AutoMigrate(&workout.Exercise{}); exercise_result.Error != nil {
		return exercise_result.Error
	}

	return nil
}
