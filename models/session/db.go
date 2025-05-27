package session

import (
	"acore/database/db"
	"fmt"

	"github.com/google/uuid"
)

type SessionID struct {
	ID uuid.UUID `db:"id"`
}

func dbGetUserSessionByID(userID uuid.UUID, token string) (uuid.UUID, error) {
	sessionID, err := db.CallFuncSingle[SessionID](db.CallFuncParams{
		FuncName: "get_user_session_by_id",
		FuncArgs: []interface{}{userID, token},
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("dbCreateUserSessionByID: %w", err)
	}
	return sessionID.ID, err
}

func dbCreateUserSession(p SessionParams) (*Session, error) {
	out, err := db.CallFuncSingle[Session](db.CallFuncParams{
		FuncName: "create_user_session",
		FuncArgs: []interface{}{
			p.UserID,
			int16(p.Type),
			int16(p.Provider),
			p.Token,
			p.IPAddress,
			p.UserAgent,
			p.ExpiresAt,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("dbCreateUserSession: %w", err)
	}
	return out, nil
}
