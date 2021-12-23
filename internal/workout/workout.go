package workout

import "github.com/jinzhu/gorm"

// Service : comment service
type Service struct {
	DB *gorm.DB
}

// Workout : struct
type Workout struct {
	gorm.Model
	Date     string
	Type     string
	Duration string
	Exercise string
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
func (s *Service) GetWorkout(ID uint) (Workout, error) {
	var workout Workout
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
func (s *Service) UpdateComment(ID uint, newWorkout Workout) (Workout, error) {
	workout, err := s.GetWorkout(ID)

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
	if _, err := s.GetWorkout(ID); err != nil {
		return err
	}

	if result := s.DB.Delete(&Workout{}, ID); result.Error != nil {
		return result.Error
	}

	return nil
}

// GetAllWorkouts : get all workouts
func (s *Service) GetAllWorkouts() ([]Workout, error) {
	var workouts []Workout
	if result := s.DB.Find(&workouts); result.Error != nil {
		return workouts, result.Error
	}
	return workouts, nil
}
