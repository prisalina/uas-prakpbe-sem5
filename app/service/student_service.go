package service

import "uas-pbe-praksem5/app/repository"

type StudentService struct {
	Repo *repository.StudentRepo
}

func NewStudentService(r *repository.StudentRepo) *StudentService { return &StudentService{Repo: r} }

func (s *StudentService) CreateStudent(req interface{}) error {
	// implement accordingly if you want typed model param
	return nil
}
