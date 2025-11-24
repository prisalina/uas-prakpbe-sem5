package repository

import (
	"database/sql"
	"uas-pbe-praksem5/app/model"
)

type StudentRepo struct {
	DB *sql.DB
}

func NewStudentRepo(db *sql.DB) *StudentRepo { return &StudentRepo{DB: db} }

func (r *StudentRepo) Create(s model.Student) error {
	_, err := r.DB.Exec(`INSERT INTO students (user_id,student_id,program_study,academic_year,advisor_id,created_at)
	VALUES ($1,$2,$3,$4,$5,NOW())`, s.UserID, s.StudentID, s.ProgramStudy, s.AcademicYear, s.AdvisorID)
	return err
}

func (r *StudentRepo) GetByUserID(userID string) (*model.Student, error) {
	row := r.DB.QueryRow(`SELECT id,user_id,student_id,program_study,academic_year,advisor_id,created_at FROM students WHERE user_id=$1`, userID)
	var s model.Student
	if err := row.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt); err != nil {
		return nil, err
	}
	return &s, nil
}
