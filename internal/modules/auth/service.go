package auth

import (
	"errors"
	"fmt"

	"github.com/cloudmachinery/movie-reviews/internal/modules/jwt"
	"github.com/cloudmachinery/movie-reviews/internal/modules/users"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type Service struct {
	userService *users.Service
	jwtService  *jwt.Service
}

func NewService(userService *users.Service, jwtService *jwt.Service) *Service {
	return &Service{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (s *Service) Register(ctx context.Context, user *users.User, password string) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userWithPassword := &users.UserWithPassword{
		User:         user,
		PasswordHash: string(passHash),
	}

	return s.userService.Create(ctx, userWithPassword)
}

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userService.GetUserWithPassword(ctx, email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("wrong email or password")
	}

	accessToken, err := s.jwtService.GenerateToken(user.Id, user.Role)
	if err != nil {
		return "", fmt.Errorf("jwt.GenerateToken: %w", err)
	}

	return accessToken, nil
}
