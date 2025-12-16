package cache

import (
    "Gym_StrongCode/internal/models"
    "fmt"
    "time"
)

type SessionManager struct {
    redisClient *RedisClient
}

func NewSessionManager(redisClient *RedisClient) *SessionManager {
    return &SessionManager{
        redisClient: redisClient,
    }
}

// CreateSession создает новую сессию
func (sm *SessionManager) CreateSession(userID int, userEmail string, isAdmin bool) (string, error) {
    sessionID := generateSessionID()
    
    session := models.Session{
        ID:        sessionID,
        UserID:    userID,
        UserEmail: userEmail,
        IsAdmin:   isAdmin,
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }

    err := sm.redisClient.Set(
        fmt.Sprintf("session:%s", sessionID),
        session,
        24*time.Hour,
    )
    if err != nil {
        return "", fmt.Errorf("failed to create session: %w", err)
    }

    // Также сохраняем сессию по userID для быстрого поиска
    err = sm.redisClient.Set(
        fmt.Sprintf("user_sessions:%d", userID),
        sessionID,
        24*time.Hour,
    )
    if err != nil {
        return "", fmt.Errorf("failed to store user session: %w", err)
    }

    return sessionID, nil
}

// GetSession получает сессию по ID
func (sm *SessionManager) GetSession(sessionID string) (*models.Session, error) {
    var session models.Session
    err := sm.redisClient.Get(
        fmt.Sprintf("session:%s", sessionID),
        &session,
    )
    if err != nil {
        return nil, fmt.Errorf("session not found: %w", err)
    }

    // Проверяем не истекла ли сессия
    if time.Now().After(session.ExpiresAt) {
        sm.DeleteSession(sessionID)
        return nil, fmt.Errorf("session expired")
    }

    return &session, nil
}

// DeleteSession удаляет сессию
func (sm *SessionManager) DeleteSession(sessionID string) error {
    // Получаем сессию чтобы узнать userID
    var session models.Session
    err := sm.redisClient.Get(
        fmt.Sprintf("session:%s", sessionID),
        &session,
    )
    if err == nil {
        // Удаляем связку userID -> sessionID
        sm.redisClient.Delete(fmt.Sprintf("user_sessions:%d", session.UserID))
    }

    // Удаляем саму сессию
    return sm.redisClient.Delete(fmt.Sprintf("session:%s", sessionID))
}

// GetUserSession получает активную сессию пользователя
func (sm *SessionManager) GetUserSession(userID int) (*models.Session, error) {
    var sessionID string
    err := sm.redisClient.Get(
        fmt.Sprintf("user_sessions:%d", userID),
        &sessionID,
    )
    if err != nil {
        return nil, fmt.Errorf("no active session found for user: %w", err)
    }

    return sm.GetSession(sessionID)
}

// ExtendSession продлевает сессию
func (sm *SessionManager) ExtendSession(sessionID string, duration time.Duration) error {
    session, err := sm.GetSession(sessionID)
    if err != nil {
        return err
    }

    session.ExpiresAt = time.Now().Add(duration)
    
    return sm.redisClient.Set(
        fmt.Sprintf("session:%s", sessionID),
        session,
        duration,
    )
}

// GetAllActiveSessions возвращает все активные сессии
func (sm *SessionManager) GetAllActiveSessions() ([]models.Session, error) {
    keys, err := sm.redisClient.Keys("session:*")
    if err != nil {
        return nil, err
    }

    var sessions []models.Session
    for _, key := range keys {
        var session models.Session
        if err := sm.redisClient.Get(key, &session); err == nil {
            if time.Now().Before(session.ExpiresAt) {
                sessions = append(sessions, session)
            }
        }
    }

    return sessions, nil
}

// CountActiveSessions подсчитывает активные сессии
func (sm *SessionManager) CountActiveSessions() (int64, error) {
    keys, err := sm.redisClient.Keys("session:*")
    if err != nil {
        return 0, err
    }

    count := int64(0)
    for _, key := range keys {
        var session models.Session
        if err := sm.redisClient.Get(key, &session); err == nil {
            if time.Now().Before(session.ExpiresAt) {
                count++
            }
        }
    }

    return count, nil
}

// ClearExpiredSessions очищает истекшие сессии
func (sm *SessionManager) ClearExpiredSessions() (int64, error) {
    keys, err := sm.redisClient.Keys("session:*")
    if err != nil {
        return 0, err
    }

    deleted := int64(0)
    for _, key := range keys {
        var session models.Session
        if err := sm.redisClient.Get(key, &session); err == nil {
            if time.Now().After(session.ExpiresAt) {
                if err := sm.redisClient.Delete(key); err == nil {
                    // Также удаляем связку userID -> sessionID
                    sm.redisClient.Delete(fmt.Sprintf("user_sessions:%d", session.UserID))
                    deleted++
                }
            }
        }
    }

    return deleted, nil
}

func generateSessionID() string {
    return fmt.Sprintf("session_%d", time.Now().UnixNano())
}