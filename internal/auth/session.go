package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// type session struct {
// 	session_id string
// 	user_id    int
// 	created_at time.Time
// 	expires_at time.Time
// }

func GenerateSession(userID int) (string, error) {
	sessionID := uuid.NewString()
	createdAt := time.Now()
	expiresAt := time.Now().Add(time.Hour * 2)
	db, err := sql.Open("sqlite3", "./internal/db/app.db")
	if err != nil {
		return "", fmt.Errorf("failed to open database connection : %v", err)
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO sessions (session_id,user_id,created_at,expires_at) VALUES (?,?,?,?)", sessionID, userID, createdAt, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to create session : %v", err)
	}
	return sessionID, nil
}

func GetUserIdFromSessions(w http.ResponseWriter, sessionID string) (int, error) {
	var userID int
	var expiresAt time.Time

	db, err := sql.Open("sqlite3", "./internal/db/app.db")
	if err != nil {
		return userID, fmt.Errorf("failed to open database connection : %v", err)
	}
	defer db.Close()

	err = db.QueryRow("SELECT user_id,expires_at FROM sessions WHERE session_id = ?", sessionID).Scan(&userID, &expiresAt)
	if err != nil {
		return userID, fmt.Errorf("failed to fetch user session from database : %v", err)
	}
	if time.Now().After(expiresAt) {
		err = DeleteSession(w, sessionID)
		if err != nil {
			return userID, fmt.Errorf("failed to remove expired session from database : %v", err)
		}
		return userID, fmt.Errorf("session has expired")
	}
	return userID, nil
}

func DeleteSession(w http.ResponseWriter, sessionID string) error {
	db, err := sql.Open("sqlite3", "./internal/db/app.db")
	if err != nil {
		return fmt.Errorf("failed to open database connection : %v", err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	if err != nil {
		return fmt.Errorf("failed to remove expired session from database : %v", err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
	return nil
}
