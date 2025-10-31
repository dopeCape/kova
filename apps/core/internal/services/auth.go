package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dopeCape/kova/internal/models"
	"github.com/dopeCape/kova/internal/store"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService *UserService
	store       store.Store
	jwtSecret   []byte
	validator   *validator.Validate
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

func NewAuthService(userService *UserService, store store.Store, jwtSecret string) *AuthService {
	return &AuthService{
		userService: userService,
		store:       store,
		jwtSecret:   []byte(jwtSecret),
		validator:   validator.New(),
	}
}

// Login authenticates a user and returns JWT tokens
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.User, *TokenPair, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get user by email or username
	user, err := s.userService.GetUserByEmailOrUsername(ctx, req.Login)
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// Generate tokens
	tokenPair, err := s.generateTokenPair(user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user.ToPublic(), tokenPair, nil
}

// RefreshToken generates new tokens using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Parse and validate refresh token
	claims, err := s.validateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user to ensure they still exist
	user, err := s.store.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new tokens
	tokenPair, err := s.generateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokenPair, nil
}

// ValidateToken validates a JWT token and returns claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	return s.validateToken(tokenString)
}

// GetCurrentUser returns the current user from token
func (s *AuthService) GetCurrentUser(ctx context.Context, userID string) (*models.User, error) {
	return s.userService.GetUser(ctx, userID)
}

// Logout invalidates tokens (in a real implementation, you'd maintain a blacklist)
func (s *AuthService) Logout(ctx context.Context, tokenString string) error {
	// In a production system, you would:
	// 1. Add token to blacklist/redis
	// 2. Or use shorter token expiry with refresh rotation
	// For now, we'll just validate the token exists
	_, err := s.validateToken(tokenString)
	if err != nil {
		return errors.New("invalid token")
	}

	// TODO: Add to token blacklist
	return nil
}

// generateTokenPair creates access and refresh tokens
func (s *AuthService) generateTokenPair(user *models.User) (*TokenPair, error) {
	now := time.Now()

	// Access token (15 minutes)
	accessClaims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(365 * 24 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "kova",
			Subject:   user.ID,
			ID:        fmt.Sprintf("access_%d", now.Unix()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshClaims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "deployment-manager",
			Subject:   user.ID,
			ID:        fmt.Sprintf("refresh_%d", now.Unix()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    15 * 60, // 15 minutes in seconds
	}, nil
}

// validateToken parses and validates a JWT token
func (s *AuthService) validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) ExtractTokenFromHeader(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	// Expected format: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// ChangePassword allows authenticated user to change password
func (s *AuthService) ChangePassword(ctx context.Context, userID string, req *models.ChangePasswordRequest) error {
	return s.userService.ChangePassword(ctx, userID, req)
}
