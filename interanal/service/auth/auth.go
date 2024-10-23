package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/interanal/domain/models"
	"sso/interanal/storage"
	ssojwt "sso/lib/jwt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCreds = errors.New("invalid creds")
	ErrAppNotFound  = errors.New("app not found")
	ErrUserExists   = errors.New("user already exists")
)

type AuthService struct {
	logger             *slog.Logger
	userSaver          UserSaver
	userProvider       UserProvider
	appServiceProvider AppServiceProvider
	tokenTTL           time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passwordHash []byte,
	) (userId int64, err error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type AppServiceProvider interface {
	GetApp(ctx context.Context, appId int) (models.App, error)
}

func New(
	logger *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appServiceProvider AppServiceProvider,
	tokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		logger:             logger,
		userSaver:          userSaver,
		userProvider:       userProvider,
		appServiceProvider: appServiceProvider,
		tokenTTL:           tokenTTL,
	}
}

func (a *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
	appId int,
) (string, error) {
	const op = "service.auth.Login"
	logger := a.logger.With(slog.String("op", op))
	logger.Info("login user", slog.Int("app_id", appId))

	user, err := a.userProvider.GetUser(ctx, email)
	if errors.Is(err, storage.ErrUserNotFound) {
		logger.Warn("user not found")
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCreds)
	}
	if err != nil {
		logger.Error("failed to get user", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PaswordHash, []byte(password)); err != nil {
		logger.Warn("invalid creds")
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCreds)
	}

	app, err := a.appServiceProvider.GetApp(ctx, appId)
	if errors.Is(err, storage.ErrAppNotFound) {
		logger.Warn("app not found")
		return "", fmt.Errorf("%s: %w", op, ErrAppNotFound)
	}
	if err != nil {
		logger.Error("failed to get app", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	logger.Info("user logged successfully")
	token, err := ssojwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.logger.Error("failed to generate token", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}

func (a *AuthService) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	const op = "service.auth.RegisterNewUser"
	logger := a.logger.With(slog.String("op", op))
	logger.Info("register user")
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to generate password hash", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	userId, err := a.userSaver.SaveUser(ctx, email, passwordHash)
	if errors.Is(err, storage.ErrUserExists) {
		logger.Warn("user already exists", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
	}
	if err != nil {
		logger.Error("failed to save user", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	logger.Info("user saved successfully")
	return userId, nil
}

func (a *AuthService) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "service.auth.IsAdmin"
	logger := a.logger.With(slog.String("op", op))
	logger.Info("checking if user id admin")
	isAdmin, err := a.userProvider.IsAdmin(ctx, userId)
	if err != nil {
		logger.Error("failed to determinate admin")
		return false, fmt.Errorf("%s: %w", op, err)
	}
	logger.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))
	return isAdmin, nil
}
