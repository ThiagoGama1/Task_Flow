package tests

import (
	"errors"
	"taskflow/internal/models"
)

// ─── Mock UserRepository ────────────────────────────────────────────────────

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

// ─── Mock ProjectRepository ─────────────────────────────────────────────────

type mockProjectRepo struct {
	projects map[uint]*models.Project
	members  map[uint]map[uint]bool // projectID → set of userIDs
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

// ─── Mock TaskRepository ────────────────────────────────────────────────────

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

func (m *mockTaskRepo) Delete(id uint) error {
	delete(m.tasks, id)
	return nil
}
