package workout

import (
	"github.com/jinzhu/gorm"
)

// Service : comment service
type Service struct {
	DB *gorm.DB
}

// Workout : struct
type Workout struct {
	gorm.Model
	Date      string
	Exercises []Exercise
}

// Exercise : exercises
type Exercise struct {
	gorm.Model
	WorkoutID uint
	Name      string
	Weight    uint
	Sets      uint
	Reps      uint
	Warmup    bool
}

// CommentService : interface for workout service
type WorkoutService interface {
	GetWorkout(ID uint) (Workout, error)
	AddWorkout(workout Workout) (Workout, error)
	UpdateWorkout(ID uint, newWorkout Workout) (Workout, error)
	DeleteWorkout(ID uint) error
	GetAllWorkouts() ([]Workout, error)
}

// NewService : new workout service
func NewService(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}

// GetWorkout : get workout
func (s *Service) GetWorkout(db *gorm.DB, ID uint) (Workout, error) {
	var workout Workout

	db.Preload("Exercises").Where("id = ? ", ID).First(&workout)

	if result := s.DB.First(&workout, ID); result.Error != nil {
		return Workout{}, result.Error
	}
	return workout, nil
}

// AddWorkout : add new workout
func (s *Service) AddWorkout(workout Workout) (Workout, error) {
	if result := s.DB.Save(&workout); result.Error != nil {
		return Workout{}, result.Error
	}
	return workout, nil
}

// UpdateWorkout : update workout
func (s *Service) UpdateWorkout(ID uint, newWorkout Workout) (Workout, error) {
	workout, err := s.GetWorkout(s.DB, ID)

	if err != nil {
		return Workout{}, err
	}

	if result := s.DB.Model(&workout).Updates(newWorkout); result.Error != nil {
		return Workout{}, result.Error
	}

	return workout, nil
}

// DeleteWorkout : delete workout
func (s *Service) DeleteWorkout(ID uint) error {
	// check if workout exists
	if _, err := s.GetWorkout(s.DB, ID); err != nil {
		return err
	}

	if result := s.DB.Delete(&Workout{}, ID); result.Error != nil {
		return result.Error
	}

	return nil
}

// GetAllWorkouts : get all workouts
func (s *Service) GetAllWorkouts(db *gorm.DB) ([]Workout, error) {
	var workouts []Workout

	if result := db.Preload("Exercises").Find(&workouts); result.Error != nil {
		return []Workout{}, result.Error
	}
	return workouts, nil
}
