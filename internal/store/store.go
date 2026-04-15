package store

import (
	"errors"
	"sync"
)

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type MemoryStore struct {
	mu    sync.RWMutex
	users map[string]User
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{users: make(map[string]User)}
}

func (s *MemoryStore) Exists(username, email, mobile string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if u.Username == username {
			return errors.New("username taken")
		}
		if u.Email == email {
			return errors.New("email registered")
		}
		if u.Mobile == mobile {
			return errors.New("mobile registered")
		}
	}
	return nil
}

func (s *MemoryStore) Save(u User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[u.Username] = u
}

func (s *MemoryStore) Authenticate(username, password string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.users[username]
	if !exists {
		return User{}, false // Return empty user and false
	}
	return user, user.Password == password // Return user and true/false
}
