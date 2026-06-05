package repositories

import (
	"taskflow/internal/models"

	"gorm.io/gorm"
)

type taskRepo struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepo{db: db}
}

func (r *taskRepo) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepo) FindByID(id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.Preload("Assignee").First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepo) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepo) FindAssignedTo(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Preload("Project").Preload("Assignee").
		Where("assignee_id = ? AND deleted_at IS NULL", userID).
		Order("due_date ASC").
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepo) Delete(id uint) error {
	return r.db.Delete(&models.Task{}, id).Error
}
