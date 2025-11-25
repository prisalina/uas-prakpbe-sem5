package service

import "uas-pbe-praksem5/app/repository"

type AchievementRefService struct {
	Repo *repository.AchievementRefRepo
}

func NewAchievementRefService(r *repository.AchievementRefRepo) *AchievementRefService {
	return &AchievementRefService{Repo: r}
}

func (s *AchievementRefService) CreateRef(studentID, mongoID string) error {
	return s.Repo.Create(studentID, mongoID)
}

func (s *AchievementRefService) UpdateStatus(id, status string, verifierID *string, rejectionNote *string) error {
	return s.Repo.UpdateStatus(id, status, verifierID, rejectionNote)
}
