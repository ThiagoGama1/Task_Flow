package repositories

import (
	"taskflow/internal/models"

	"gorm.io/gorm"
)

type commentRepo struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepo{db: db}
}

func (r *commentRepo) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepo) FindByTaskID(taskID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.
		Where("task_id = ?", taskID).
		Preload("User").
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}
