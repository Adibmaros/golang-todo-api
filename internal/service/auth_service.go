package service

import (
	"errors"
	"strings"

	"github.com/Adibmaros/todo-api/internal/auth"
	"github.com/Adibmaros/todo-api/internal/model"
	"github.com/Adibmaros/todo-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyUsed   = errors.New("email sudah terdaftar")
	ErrInvalidCredentials = errors.New("email atau password salah")
)

type AuthService interface {
	Register(email, password string) (*model.User, error)
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(email, password string) (*model.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return nil, ErrEmailAlreadyUsed
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    email,
		Password: string(hashed),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	user, err := s.userRepo.FindByEmail(email)
	if errors.Is(err, repository.ErrNotFound) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return auth.GenerateToken(user.ID)
}
