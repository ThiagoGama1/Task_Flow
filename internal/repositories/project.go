package repositories

import (
	"taskflow/internal/models"

	"gorm.io/gorm"
)

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepo{db: db}
}

func (r *projectRepo) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepo) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	err := r.db.First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepo) WithMembers(id uint) (*models.Project, error) {
	var project models.Project
	err := r.db.
		Preload("Owner").
		Preload("Members").
		Preload("Tasks").
		Preload("Tasks.Assignee").
		First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepo) FindByMemberID(userID uint) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.
		Joins("JOIN project_members ON project_members.project_id = projects.id").
		Where("project_members.user_id = ?", userID).
		Preload("Owner").
		Order("projects.created_at DESC").
		Find(&projects).Error
	return projects, err
}

func (r *projectRepo) AddMember(projectID, userID uint) error {
	project := models.Project{}
	project.ID = projectID
	user := models.User{}
	user.ID = userID
	return r.db.Model(&project).Association("Members").Append(&user)
}

func (r *projectRepo) RemoveMember(projectID, userID uint) error {
	project := models.Project{}
	project.ID = projectID
	user := models.User{}
	user.ID = userID
	return r.db.Model(&project).Association("Members").Delete(&user)
}

func (r *projectRepo) IsMember(projectID, userID uint) bool {
	var count int64
	r.db.Table("project_members").
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count)
	return count > 0
}

func (r *projectRepo) Delete(id uint) error {
	return r.db.Delete(&models.Project{}, id).Error
}
