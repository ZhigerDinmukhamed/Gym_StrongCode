package cache

import (
	"Gym_StrongCode/internal/models"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

const sessionTTL = 24 * time.Hour

type SessionManager struct {
	redisClient *RedisClient
}

func NewSessionManager(redisClient *RedisClient) *SessionManager {
	return &SessionManager{
		redisClient: redisClient,
	}
}

// CreateSession создает новую сессию (1 пользователь = 1 сессия)
func (sm *SessionManager) CreateSession(
	userID int,
	userEmail string,
	isAdmin bool,
) (string, error) {

	sessionID, err := generateSessionID()
	if err != nil {
		return "", err
	}

	session := models.Session{
		ID:        sessionID,
		UserID:    userID,
		UserEmail: userEmail,
		IsAdmin:   isAdmin,
		CreatedAt: time.Now(),
	}

	sessionKey := fmt.Sprintf("session:%s", sessionID)
	userKey := fmt.Sprintf("user_session:%d", userID)

	// сохраняем сессию
	if err := sm.redisClient.Set(sessionKey, session, sessionTTL); err != nil {
		return "", err
	}

	// привязка user -> session
	if err := sm.redisClient.Set(userKey, sessionID, sessionTTL); err != nil {
		_ = sm.redisClient.Delete(sessionKey)
		return "", err
	}

	return sessionID, nil
}

// GetSession возвращает сессию по sessionID
func (sm *SessionManager) GetSession(sessionID string) (*models.Session, error) {
	var session models.Session

	err := sm.redisClient.Get(
		fmt.Sprintf("session:%s", sessionID),
		&session,
	)
	if err != nil {
		return nil, fmt.Errorf("session not found")
	}

	return &session, nil
}

// GetUserSession возвращает активную сессию пользователя
func (sm *SessionManager) GetUserSession(userID int) (*models.Session, error) {
	var sessionID string

	err := sm.redisClient.Get(
		fmt.Sprintf("user_session:%d", userID),
		&sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("no active session")
	}

	return sm.GetSession(sessionID)
}

// ExtendSession продлевает сессию
func (sm *SessionManager) ExtendSession(sessionID string) error {
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return err
	}

	sessionKey := fmt.Sprintf("session:%s", sessionID)
	userKey := fmt.Sprintf("user_session:%d", session.UserID)

	if err := sm.redisClient.Set(sessionKey, session, sessionTTL); err != nil {
		return err
	}

	if err := sm.redisClient.Set(userKey, sessionID, sessionTTL); err != nil {
		return err
	}

	return nil
}

// DeleteSession удаляет сессию
func (sm *SessionManager) DeleteSession(sessionID string) error {
	session, err := sm.GetSession(sessionID)
	if err == nil {
		_ = sm.redisClient.Delete(
			fmt.Sprintf("user_session:%d", session.UserID),
		)
	}

	return sm.redisClient.Delete(fmt.Sprintf("session:%s", sessionID))
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
