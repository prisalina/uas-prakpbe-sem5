package repository

import (
	"database/sql"
)

type AchievementRefRepo struct {
	DB *sql.DB
}

func NewAchievementRefRepo(db *sql.DB) *AchievementRefRepo { return &AchievementRefRepo{DB: db} }

func (r *AchievementRefRepo) Create(studentID, mongoID string) error {
	_, err := r.DB.Exec(`INSERT INTO achievement_references (student_id,mongo_achievement_id,status,created_at,updated_at)
	VALUES ($1,$2,'draft',NOW(),NOW())`, studentID, mongoID)
	return err
}

func (r *AchievementRefRepo) UpdateStatus(id, status string, verifierID *string, rejectionNote *string) error {
	if verifierID != nil {
		_, err := r.DB.Exec(`UPDATE achievement_references SET status=$1, verified_at=NOW(), verified_by=$2, updated_at=NOW() WHERE id=$3`, status, *verifierID, id)
		return err
	}
	if rejectionNote != nil {
		_, err := r.DB.Exec(`UPDATE achievement_references SET status=$1, rejection_note=$2, updated_at=NOW() WHERE id=$3`, status, *rejectionNote, id)
		return err
	}
	_, err := r.DB.Exec(`UPDATE achievement_references SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}
