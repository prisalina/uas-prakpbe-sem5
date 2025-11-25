package service

import (
	"uas-pbe-praksem5/app/model"
	"uas-pbe-praksem5/app/repository"
)

type LecturerService struct {
	Repo *repository.LecturerRepo
}

func NewLecturerService(r *repository.LecturerRepo) *LecturerService {
	return &LecturerService{Repo: r}
}

func (s *LecturerService) CreateLecturer(req model.Lecturer) error {
	return s.Repo.Create(req)
}
