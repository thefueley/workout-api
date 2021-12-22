package comment

import "github.com/jinzhu/gorm"

// Service : comment service
type Service struct {
	DB *gorm.DB
}

// Comment : struct
type Comment struct {
	gorm.Model
	Slug   string
	Body   string
	Author string
}

// CommentService : interface for comment service
type CommentService interface {
	GetComment(ID uint) (Comment, error)
	GetCommentBySlug(slug string) ([]Comment, error)
	PostComment(comment Comment) (Comment, error)
	UpdateComment(ID uint, newComment Comment) (Comment, error)
	DeleteComment(ID uint) error
	GetAllComments() ([]Comment, error)
}

// NewService : return new comment service
func NewService(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}

// GetComment : get comment by id
func (s *Service) GetComment(ID uint) (Comment, error) {
	var comment Comment
	if result := s.DB.First(&comment, ID); result.Error != nil {
		return Comment{}, result.Error
	}
	return comment, nil
}

// GetCommentBySlug : retrieves all comments by slug (/path)
func (s *Service) GetCommentBySlug(slug string) ([]Comment, error) {
	var comments []Comment
	if result := s.DB.Find(&comments).Where("slug = ?", slug); result.Error != nil {
		return []Comment{}, result.Error
	}
	return comments, nil
}

// PostComment : add new comment to db
func (s *Service) PostComment(comment Comment) (Comment, error) {
	if result := s.DB.Save(&comment); result.Error != nil {
		return Comment{}, result.Error
	}
	return comment, nil
}

// UpdateComment : update comment by id in db
func (s *Service) UpdateComment(ID uint, newComment Comment) (Comment, error) {
	comment, err := s.GetComment(ID)
	if err != nil {
		return Comment{}, err
	}

	if result := s.DB.Model(&comment).Updates(newComment); result.Error != nil {
		return Comment{}, result.Error
	}

	return comment, nil
}

// DeleteComment : delete comment by id from db
func (s *Service) DeleteComment(ID uint) error {
	// check if comment id exists
	if _, err := s.GetComment(ID); err != nil {
		return err
	}

	if result := s.DB.Delete(&Comment{}, ID); result.Error != nil {
		return result.Error
	}

	return nil
}

// GetAllComments : get all comments from DB
func (s *Service) GetAllComments() ([]Comment, error) {
	var comments []Comment
	if result := s.DB.Find(&comments); result.Error != nil {
		return comments, result.Error
	}
	return comments, nil
}
