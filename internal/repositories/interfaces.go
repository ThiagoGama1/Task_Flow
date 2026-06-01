package repositories

import "taskflow/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
}

type ProjectRepository interface {
	Create(project *models.Project) error
	FindByID(id uint) (*models.Project, error)
	WithMembers(id uint) (*models.Project, error)
	FindByMemberID(userID uint) ([]models.Project, error)
	AddMember(projectID, userID uint) error
	RemoveMember(projectID, userID uint) error
	IsMember(projectID, userID uint) bool
	Delete(id uint) error
}

type TaskRepository interface {
	Create(task *models.Task) error
	FindByID(id uint) (*models.Task, error)
	Update(task *models.Task) error
	Delete(id uint) error
}
