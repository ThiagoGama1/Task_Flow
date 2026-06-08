package tests

import (
	"errors"
	"taskflow/internal/models"
	"time"
)

type mockUserRepo struct {
	users  map[uint]*models.User
	nextID uint
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[uint]*models.User), nextID: 1}
}

func (m *mockUserRepo) Create(user *models.User) error {
	for _, u := range m.users {
		if u.Email == user.Email {
			return errors.New("email already exists")
		}
	}
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) FindByEmail(email string) (*models.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) FindByID(id uint) (*models.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) Update(u *models.User) error {
	m.users[u.ID] = u
	return nil
}

type mockCommentRepo struct {
	comments map[uint]*models.Comment
	nextID   uint
}

func newMockCommentRepo() *mockCommentRepo {
	return &mockCommentRepo{comments: make(map[uint]*models.Comment), nextID: 1}
}

func (m *mockCommentRepo) Create(c *models.Comment) error {
	c.ID = m.nextID
	m.nextID++
	m.comments[c.ID] = c
	return nil
}

func (m *mockCommentRepo) FindByTaskID(taskID uint) ([]models.Comment, error) {
	var out []models.Comment
	for _, c := range m.comments {
		if c.TaskID == taskID {
			out = append(out, *c)
		}
	}
	return out, nil
}

type mockProjectRepo struct {
	projects map[uint]*models.Project
	members  map[uint]map[uint]bool
	nextID   uint
}

func newMockProjectRepo() *mockProjectRepo {
	return &mockProjectRepo{
		projects: make(map[uint]*models.Project),
		members:  make(map[uint]map[uint]bool),
		nextID:   1,
	}
}

func (m *mockProjectRepo) Create(p *models.Project) error {
	p.ID = m.nextID
	m.nextID++
	m.projects[p.ID] = p
	m.members[p.ID] = make(map[uint]bool)
	return nil
}

func (m *mockProjectRepo) FindByID(id uint) (*models.Project, error) {
	if p, ok := m.projects[id]; ok {
		return p, nil
	}
	return nil, errors.New("not found")
}

func (m *mockProjectRepo) WithMembers(id uint) (*models.Project, error) {
	return m.FindByID(id)
}

func (m *mockProjectRepo) FindByMemberID(userID uint) ([]models.Project, error) {
	var out []models.Project
	for id, members := range m.members {
		if members[userID] {
			if p, ok := m.projects[id]; ok {
				out = append(out, *p)
			}
		}
	}
	return out, nil
}

func (m *mockProjectRepo) AddMember(projectID, userID uint) error {
	if _, ok := m.members[projectID]; !ok {
		m.members[projectID] = make(map[uint]bool)
	}
	m.members[projectID][userID] = true
	return nil
}

func (m *mockProjectRepo) RemoveMember(projectID, userID uint) error {
	if m.members[projectID] != nil {
		delete(m.members[projectID], userID)
	}
	return nil
}

func (m *mockProjectRepo) IsMember(projectID, userID uint) bool {
	return m.members[projectID][userID]
}

func (m *mockProjectRepo) Delete(id uint) error {
	delete(m.projects, id)
	delete(m.members, id)
	return nil
}

type mockTaskRepo struct {
	tasks  map[uint]*models.Task
	nextID uint
}

func newMockTaskRepo() *mockTaskRepo {
	return &mockTaskRepo{tasks: make(map[uint]*models.Task), nextID: 1}
}

func (m *mockTaskRepo) Create(t *models.Task) error {
	t.ID = m.nextID
	m.nextID++
	m.tasks[t.ID] = t
	return nil
}

func (m *mockTaskRepo) FindByID(id uint) (*models.Task, error) {
	if t, ok := m.tasks[id]; ok {
		return t, nil
	}
	return nil, errors.New("not found")
}

func (m *mockTaskRepo) Update(t *models.Task) error {
	m.tasks[t.ID] = t
	return nil
}

func (m *mockTaskRepo) FindAssignedTo(userID uint) ([]models.Task, error) {
	var out []models.Task
	for _, t := range m.tasks {
		if t.AssigneeID != nil && *t.AssigneeID == userID {
			out = append(out, *t)
		}
	}
	return out, nil
}

func (m *mockTaskRepo) Delete(id uint) error {
	delete(m.tasks, id)
	return nil
}

type mockActivityRepo struct {
	logs   []*models.ActivityLog
	nextID uint
}

func newMockActivityRepo() *mockActivityRepo {
	return &mockActivityRepo{nextID: 1}
}

func (m *mockActivityRepo) Create(log *models.ActivityLog) error {
	log.ID = m.nextID
	m.nextID++
	m.logs = append(m.logs, log)
	return nil
}

func (m *mockActivityRepo) FindByProjectID(projectID uint, limit int) ([]models.ActivityLog, error) {
	var out []models.ActivityLog
	for _, l := range m.logs {
		if l.ProjectID == projectID {
			out = append(out, *l)
		}
	}
	if len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (m *mockActivityRepo) CountForUser(userID uint, since time.Time) (int64, error) {
	return 0, nil
}
