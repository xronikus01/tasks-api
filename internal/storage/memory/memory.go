package memory

import (
	"errors"
	"sort"
	"sync"
	"time"

	"tasks-api/internal/models"
)

var (
	ErrNotFound = errors.New("task not found")
)

type Store struct {
	mu     sync.RWMutex
	nextID int
	tasks  map[int]models.Task
}

func New() *Store {
	return &Store{
		nextID: 1,
		tasks:  make(map[int]models.Task),
	}
}

func (s *Store) List() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		out = append(out, t)
	}

	// стабильный порядок (удобно при проверке)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (s *Store) Create(t models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t.ID = s.nextID
	s.nextID++

	if t.CreatedAt == "" {
		t.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	s.tasks[t.ID] = t
	return t, nil
}

func (s *Store) Get(id int) (models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tasks[id]
	return t, ok
}

func (s *Store) Update(id int, t models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	old, ok := s.tasks[id]
	if !ok {
		return models.Task{}, ErrNotFound
	}

	// сервер управляет ID и created_at
	t.ID = id
	t.CreatedAt = old.CreatedAt

	s.tasks[id] = t
	return t, nil
}

func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.tasks, id)
	return nil
}
