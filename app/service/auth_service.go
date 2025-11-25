package service

import (
	"errors"
	"time"
	"uas-pbe-praksem5/app/model"
	"uas-pbe-praksem5/app/repository"
	"uas-pbe-praksem5/utils"
)

type AuthService struct {
	UserRepo         *repository.UserRepo
	RefreshTokenRepo *repository.RefreshTokenRepo
}

func NewAuthService(u *repository.UserRepo, r *repository.RefreshTokenRepo) *AuthService {
	return &AuthService{UserRepo: u, RefreshTokenRepo: r}
}

func (s *AuthService) Login(req model.LoginRequest) (string, string, *model.User, error) {
	user, pwdHash, err := s.UserRepo.GetByUsernameOrEmail(req.Username)
	if err != nil {
		return "", "", nil, err
	}
	if user == nil {
		return "", "", nil, nil
	}
	if !utils.CheckPassword(pwdHash, req.Password) {
		return "", "", nil, nil
	}
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.RoleName)
	if err != nil {
		return "", "", nil, err
	}
	refreshToken, jti, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", nil, err
	}
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := s.RefreshTokenRepo.Create(jti, user.ID, expiresAt); err != nil {
		return "", "", nil, err
	}
	return accessToken, refreshToken, user, nil
}

func (s *AuthService) Refresh(oldRefreshToken string) (string, string, error) {
	claims, err := utils.ValidateRefreshToken(oldRefreshToken)
	if err != nil {
		return "", "", err
	}
	ok, err := s.RefreshTokenRepo.ExistsAndNotExpired(claims.JTI)
	if err != nil {
		return "", "", err
	}
	if !ok {
		return "", "", errors.New("refresh token not found or expired")
	}
	if err := s.RefreshTokenRepo.DeleteByJTI(claims.JTI); err != nil {
		return "", "", err
	}
	user, err := s.UserRepo.GetByID(claims.UserID)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", errors.New("user not found")
	}
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.RoleName)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, newJTI, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}
	newExpires := time.Now().Add(7 * 24 * time.Hour)
	if err := s.RefreshTokenRepo.Create(newJTI, user.ID, newExpires); err != nil {
		return "", "", err
	}
	return accessToken, newRefreshToken, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return err
	}
	return s.RefreshTokenRepo.DeleteByJTI(claims.JTI)
}
