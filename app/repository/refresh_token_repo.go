package repository

import (
	"database/sql"
	"time"
)

type RefreshTokenRepo struct {
	DB *sql.DB
}

func NewRefreshTokenRepo(db *sql.DB) *RefreshTokenRepo { return &RefreshTokenRepo{DB: db} }

func (r *RefreshTokenRepo) Create(jti string, userID string, expiresAt time.Time) error {
	_, err := r.DB.Exec(`INSERT INTO refresh_tokens (jti,user_id,expires_at,created_at) VALUES ($1,$2,$3,NOW())`, jti, userID, expiresAt)
	return err
}

func (r *RefreshTokenRepo) DeleteByJTI(jti string) error {
	_, err := r.DB.Exec(`DELETE FROM refresh_tokens WHERE jti=$1`, jti)
	return err
}

func (r *RefreshTokenRepo) ExistsAndNotExpired(jti string) (bool, error) {
	row := r.DB.QueryRow(`SELECT expires_at FROM refresh_tokens WHERE jti=$1`, jti)
	var exp time.Time
	if err := row.Scan(&exp); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if time.Now().After(exp) {
		_, _ = r.DB.Exec(`DELETE FROM refresh_tokens WHERE jti=$1`, jti)
		return false, nil
	}
	return true, nil
}
