package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/SavanRajyaguru/ecommerce-go-user-service/internal/auth"
	"github.com/SavanRajyaguru/ecommerce-go-user-service/internal/cache"
	"github.com/SavanRajyaguru/ecommerce-go-user-service/models"
	"github.com/SavanRajyaguru/ecommerce-go-user-service/pkg/logger"
	"github.com/SavanRajyaguru/ecommerce-go-user-service/pkg/utils"
	"github.com/SavanRajyaguru/ecommerce-go-user-service/repository"
	"go.uber.org/zap"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo: repository.NewUserRepository(),
	}
}

func (s *UserService) Register(name, email, password string) error {
	// Check if user exists
	existingUser, _ := s.repo.FindByEmail(email)
	if existingUser != nil && existingUser.ID != 0 {
		return errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		logger.Log.Error("Failed to hash password", zap.Error(err))
		return err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	return s.repo.CreateUser(user)
}

func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) GetUserProfile(userID uint) (*models.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d", userID)

	// Try cache
	val, err := cache.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		var user models.User
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			return &user, nil
		}
	}

	// Fetch from DB
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Set cache
	userJson, _ := json.Marshal(user)
	cache.RDB.Set(ctx, cacheKey, userJson, 10*time.Minute)

	return user, nil
}
