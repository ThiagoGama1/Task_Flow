package repositories

import (
	"taskflow/internal/models"
	"time"

	"gorm.io/gorm"
)

type activityRepo struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepo{db: db}
}

func (r *activityRepo) Create(log *models.ActivityLog) error {
	return r.db.Create(log).Error
}

func (r *activityRepo) FindByProjectID(projectID uint, limit int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := r.db.
		Where("project_id = ?", projectID).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

func (r *activityRepo) CountForUser(userID uint, since time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.ActivityLog{}).
		Joins("JOIN project_members ON project_members.project_id = activity_logs.project_id").
		Where("project_members.user_id = ? AND activity_logs.user_id != ? AND activity_logs.created_at > ?",
			userID, userID, since).
		Count(&count).Error
	return count, err
}
