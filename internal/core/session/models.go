package session

import "github.com/google/uuid"

// Session model store user session
type Session struct {
	SessionID string    `json:"session_id" redis:"session_id"`
	UserID    uuid.UUID `json:"user_id" redis:"user_id"`
}
