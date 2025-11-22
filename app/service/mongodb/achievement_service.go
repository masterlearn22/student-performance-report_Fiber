package service

import (
	repo "student-performance-report/app/repository/mongodb"
)

type AchievementService struct {
	achievementRepo repo.AchievementRepository
}