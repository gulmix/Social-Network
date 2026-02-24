package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gulmix/Social-Network/internal/config"
	"github.com/gulmix/Social-Network/internal/models"
	"github.com/gulmix/Social-Network/internal/repository"
	"github.com/gulmix/Social-Network/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(repo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: repo,
		cfg:      cfg,
	}
}

func (s *AuthService) Register(email, username, password string) (*models.User, string, error) {
	exists, err := s.userRepo.EmailExists(email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, "", errors.New("email already exists")
	}

	exists, err = s.userRepo.UsernameExists(username)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return nil, "", errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:         email,
		Username:      username,
		PasswordHash:  sql.NullString{String: string(hashedPassword), Valid: true},
		EmailVerified: false,
	}

	if err = s.userRepo.CreateUser(user); err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	token, err := utils.GenerateToken(user.ID, user.Email, s.cfg)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !user.PasswordHash.Valid {
		return nil, "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, s.cfg)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *AuthService) OAuthLogin(ctx context.Context, provider string, code string) (*models.User, string, error) {
	var oauthConfig *oauth2.Config
	var getUserInfo func(context.Context, *oauth2.Token) (*utils.OAuthUserInfo, error)

	switch provider {
	case "google":
		oauthConfig = utils.GetGoogleOAuthConfig(s.cfg)
		getUserInfo = utils.GetGoogleUserInfo
	case "github":
		oauthConfig = utils.GetGitHubOAuthConfig(s.cfg)
		getUserInfo = utils.GetGitHubUserInfo
	default:
		return nil, "", errors.New("unsupported OAuth provider")
	}

	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange token: %w", err)
	}

	oauthUserInfo, err := getUserInfo(ctx, token)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user info: %w", err)
	}

	user, err := s.userRepo.GetUserByOAuthProvider(provider, oauthUserInfo.ID)
	if err == nil {
		user.AvatarURL = sql.NullString{String: oauthUserInfo.Avatar, Valid: true}
		if oauthUserInfo.Name != "" {
			user.FirstName = sql.NullString{String: oauthUserInfo.Name, Valid: true}
		}
		s.userRepo.UpdateUser(user)

		jwtToken, err := utils.GenerateToken(user.ID, user.Email, s.cfg)
		if err != nil {
			return nil, "", fmt.Errorf("failed to generate token: %w", err)
		}
		return user, jwtToken, nil
	}

	_, err = s.userRepo.GetUserByEmail(oauthUserInfo.Email)
	if err == nil {
		return nil, "", errors.New("email already registered with different account")
	}

	user = &models.User{
		Email:           oauthUserInfo.Email,
		Username:        oauthUserInfo.Email,
		EmailVerified:   true,
		OAuthProvider:   sql.NullString{String: provider, Valid: true},
		OAuthProviderID: sql.NullString{String: oauthUserInfo.ID, Valid: true},
		AvatarURL:       sql.NullString{String: oauthUserInfo.Avatar, Valid: true},
	}

	if oauthUserInfo.Name != "" {
		user.FirstName = sql.NullString{String: oauthUserInfo.Name, Valid: true}
	}

	if err = s.userRepo.CreateUser(user); err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	jwtToken, err := utils.GenerateToken(user.ID, user.Email, s.cfg)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, jwtToken, nil
}

func (s *AuthService) GetOAuthURL(provider string) (string, error) {
	var oauthConfig *oauth2.Config

	switch provider {
	case "google":
		oauthConfig = utils.GetGoogleOAuthConfig(s.cfg)
	case "github":
		oauthConfig = utils.GetGitHubOAuthConfig(s.cfg)
	default:
		return "", errors.New("unsupported OAuth provider")
	}

	return oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline), nil
}

func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	claims, err := utils.ValidateToken(tokenString, s.cfg)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
