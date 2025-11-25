package repository

import (
	"database/sql"
	"uas-pbe-praksem5/app/model"
)

type LecturerRepo struct {
	DB *sql.DB
}

func NewLecturerRepo(db *sql.DB) *LecturerRepo { return &LecturerRepo{DB: db} }

func (r *LecturerRepo) Create(l model.Lecturer) error {
	_, err := r.DB.Exec(`INSERT INTO lecturers (user_id,lecturer_id,department,created_at) VALUES ($1,$2,$3,NOW())`,
		l.UserID, l.LecturerID, l.Department)
	return err
}
