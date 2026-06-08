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

	if len(project.Tasks) > 0 {
		type countRow struct {
			TaskID uint
			Count  int
		}
		var rows []countRow
		taskIDs := make([]uint, len(project.Tasks))
		for i, t := range project.Tasks {
			taskIDs[i] = t.ID
		}
		r.db.Model(&models.Comment{}).
			Select("task_id, count(*) as count").
			Where("task_id IN ? AND deleted_at IS NULL", taskIDs).
			Group("task_id").
			Scan(&rows)
		counts := make(map[uint]int, len(rows))
		for _, row := range rows {
			counts[row.TaskID] = row.Count
		}
		for i := range project.Tasks {
			project.Tasks[i].CommentCount = counts[project.Tasks[i].ID]
		}
	}

	return &project, nil
}

func (r *projectRepo) FindByMemberID(userID uint) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.
		Joins("JOIN project_members ON project_members.project_id = projects.id").
		Where("project_members.user_id = ?", userID).
		Preload("Owner").
		Preload("Tasks").
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
	r.db.Where("project_id = ?", id).Delete(&models.Task{})
	return r.db.Delete(&models.Project{}, id).Error
}
