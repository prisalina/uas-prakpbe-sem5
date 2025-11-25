package service

import (
	"database/sql"
	"uas-pbe-praksem5/app/model"
	"uas-pbe-praksem5/app/repository"
	"uas-pbe-praksem5/utils"
)

type UserService struct {
	Repo *repository.UserRepo
	DB   *sql.DB
}

func NewUserService(r *repository.UserRepo, db *sql.DB) *UserService {
	return &UserService{Repo: r, DB: db}
}

func (s *UserService) CreateUser(req model.CreateUserRequest) error {
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}
	return s.Repo.Create(req, hash)
}

func (s *UserService) ListUsers() ([]model.User, error) {
	return s.Repo.GetAll()
}

func (s *UserService) GetByID(id string) (*model.User, error) {
	return s.Repo.GetByID(id)
}

func (s *UserService) UpdateUser(id string, req model.CreateUserRequest) error {
	return s.Repo.Update(id, req)
}

func (s *UserService) DeleteUser(id string) error {
	return s.Repo.Delete(id)
}

func (s *UserService) UpdateRole(id, roleID string) error {
	return s.Repo.UpdateRole(id, roleID)
}
