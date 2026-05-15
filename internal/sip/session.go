package sip

import (
	"encoding/json"
	"time"

	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/redis"

	"go.uber.org/zap"
)

// SessionManager manages SIP INVITE sessions (Redis-backed)
type SessionManager struct {
	log *zap.Logger
}

func NewSessionManager(log *zap.Logger) *SessionManager {
	return &SessionManager{log: log}
}

// Save saves a stream session to Redis
func (m *SessionManager) Save(stream string, session *model.SsrcTransaction) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	key := redis.InviteSession + stream
	return redis.Client.Set(redis.Ctx, key, data, 24*time.Hour).Err()
}

// Get retrieves a session by stream
func (m *SessionManager) Get(stream string) (*model.SsrcTransaction, error) {
	key := redis.InviteSession + stream
	data, err := redis.Client.Get(redis.Ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var session model.SsrcTransaction
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// Delete removes a session
func (m *SessionManager) Delete(stream string) error {
	key := redis.InviteSession + stream
	return redis.Client.Del(redis.Ctx, key).Err()
}

// SaveByCallID saves a session keyed by callId
func (m *SessionManager) SaveByCallID(callID string, session *model.SsrcTransaction) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	key := redis.InviteSession + callID
	return redis.Client.Set(redis.Ctx, key, data, 24*time.Hour).Err()
}

// GetByCallID retrieves a session by callId
func (m *SessionManager) GetByCallID(callID string) (*model.SsrcTransaction, error) {
	key := redis.InviteSession + callID
	data, err := redis.Client.Get(redis.Ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var session model.SsrcTransaction
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// DeleteByCallID removes a session by callId
func (m *SessionManager) DeleteByCallID(callID string) error {
	key := redis.InviteSession + callID
	return redis.Client.Del(redis.Ctx, key).Err()
}

// Exists checks if a session exists
func (m *SessionManager) Exists(stream string) bool {
	key := redis.InviteSession + stream
	exists, _ := redis.Client.Exists(redis.Ctx, key).Result()
	return exists > 0
}

// ListStreams returns all active stream keys
func (m *SessionManager) ListStreams(pattern string) ([]string, error) {
	key := redis.InviteSession + pattern
	return redis.Client.Keys(redis.Ctx, key).Result()
}
