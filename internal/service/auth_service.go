package service

import (
	"errors"
	"time"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage"
)

type AuthService interface {
	Register(email, password string) error
	Login(email, password string) (Tokens, error)
	Refresh(refreshToken string) (Tokens, error)
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type authService struct {
	userRepo   storage.UserRepository
	hasher     PasswordHasher
	jwt        JWTManager
	generateID func() string
}

type JWTManager interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ParseRefreshToken(token string) (string, error)
	ParseAccessToken(token string) (string, error)
}

func NewAuthService(
	userRepo storage.UserRepository,
	hasher PasswordHasher,
	jwt JWTManager,
	generateID func() string,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		hasher:     hasher,
		jwt:        jwt,
		generateID: generateID,
	}
}

func (s *authService) Register(email, password string) error {
	if email == "" || password == "" {
		return domain.ErrInvalidInput
	}

	_, err := s.userRepo.GetByEmail(email)
	if err == nil {
		return domain.ErrUserAlreadyExists
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return err
	}

	user := domain.User{
		ID:           s.generateID(),
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}

	return s.userRepo.Create(user)
}

func (s *authService) Login(email, password string) (Tokens, error) {
	if email == "" || password == "" {
		return Tokens{}, domain.ErrInvalidInput
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return Tokens{}, domain.ErrInvalidCredentials
		}
		return Tokens{}, err
	}

	// bcrypt check
	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		return Tokens{}, domain.ErrInvalidCredentials
	}

	accessToken, err := s.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return Tokens{}, err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Refresh(refreshToken string) (Tokens, error) {
	if refreshToken == "" {
		return Tokens{}, domain.ErrInvalidInput
	}

	userID, err := s.jwt.ParseRefreshToken(refreshToken)
	if err != nil {
		return Tokens{}, domain.ErrInvalidCredentials
	}

	_, err = s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return Tokens{}, domain.ErrInvalidCredentials
		}
		return Tokens{}, err
	}

	accessToken, err := s.jwt.GenerateAccessToken(userID)
	if err != nil {
		return Tokens{}, err
	}

	newRefreshToken, err := s.jwt.GenerateRefreshToken(userID)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil

}
