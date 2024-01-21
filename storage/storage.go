package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/johnnyluo/tss-director/model"
	"github.com/patrickmn/go-cache"
)

var ErrNotFound = errors.New("not found")

// Storage is an interface that defines the methods to be implemented by a storage.
type Storage interface {
	SetSession(sessionID string, participants []string) error
	GetSession(sessionID string) ([]string, error)
	DeleteSession(sessionID string) error
	GetMessage(sessionID, participantID string) ([]model.Message, error)
	SetMessage(sessionID, participantID string, message model.Message) error
	DeleteMessage(sessionID, participantID string) error
}

type InMemoryStorage struct {
	cache *cache.Cache
}

// NewInMemoryStorage returns a new in-memory storage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		cache: cache.New(time.Minute*10, time.Minute*15),
	}
}

// SetSession sets a session with a list of participants.
func (s *InMemoryStorage) SetSession(sessionID string, participants []string) error {
	return s.cache.Add(sessionID, participants, cache.DefaultExpiration)
}

// GetSession gets a session with a list of participants.
func (s *InMemoryStorage) GetSession(sessionID string) ([]string, error) {
	if participants, ok := s.cache.Get(sessionID); ok {
		return participants.([]string), nil
	}
	return nil, ErrNotFound
}

// DeleteSession deletes a session.
func (s *InMemoryStorage) DeleteSession(sessionID string) error {
	s.cache.Delete(sessionID)
	return nil
}

// GetMessage gets a message from a session and a participant.
func (s *InMemoryStorage) GetMessage(sessionID, participantID string) ([]model.Message, error) {
	key := fmt.Sprintf("%s-%s", sessionID, participantID)
	if messages, ok := s.cache.Get(key); ok {
		return messages.([]model.Message), nil
	}
	return nil, ErrNotFound
}

// SetMessage sets a message to a session and a participant.
func (s *InMemoryStorage) SetMessage(sessionID, participantID string, message model.Message) error {
	key := fmt.Sprintf("%s-%s", sessionID, participantID)
	if messages, ok := s.cache.Get(key); ok {
		messages = append(messages.([]model.Message), message)
		s.cache.Set(key, messages, cache.DefaultExpiration)
		return nil
	}
	s.cache.Set(key, []model.Message{message}, cache.DefaultExpiration)
	return nil
}

// DeleteMessage deletes a message from a session and a participant.
func (s *InMemoryStorage) DeleteMessage(sessionID, participantID string) error {
	key := fmt.Sprintf("%s-%s", sessionID, participantID)
	s.cache.Delete(key)
	return nil
}
