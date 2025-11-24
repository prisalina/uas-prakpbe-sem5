package service

import (
	"uas-pbe-praksem5/app/model"
	"uas-pbe-praksem5/app/repository"
	"uas-pbe-praksem5/utils"
)

type AuthService struct {
	UserRepo *repository.UserRepo
}

func NewAuthService(ur *repository.UserRepo) *AuthService {
	return &AuthService{UserRepo: ur}
}

func (s *AuthService) Login(req model.LoginRequest) (string, *model.User, error) {
	user, pwdHash, err := s.UserRepo.GetByUsernameOrEmail(req.Username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, nil
	}
	if !utils.CheckPassword(pwdHash, req.Password) {
		return "", nil, nil
	}
	token, err := utils.GenerateToken(user.ID, user.Username, user.RoleName)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}
