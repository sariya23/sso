package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sso/interanal/config"
	"sso/interanal/domain/models"
	"sso/interanal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	path := "../../../config/db.yaml"
	cfg := config.MustLoadDBConfig(path)
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", cfg.User, cfg.Password, cfg.Port, cfg.DBName)
	ctx := context.Background()
	s := MustNewConnection(ctx, dbURL)
	s.truncateUsers(ctx)
	os.Exit(exitVal)
}

// TestSuccessfullyConnectToDb проверяет,
// что подключение к базе данных происходит успешно.
func TestSuccessfullyConnectToDb(t *testing.T) {
	path := "../../../config/db.yaml"
	cfg := config.MustLoadDBConfig(path)
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", cfg.User, cfg.Password, cfg.Port, cfg.DBName)
	ctx := context.Background()

	MustNewConnection(ctx, dbURL)
}

// TestSaveUser проверяет, что пользователь
// успешно сохраняется.
func TestSaveUser(t *testing.T) {
	path := "../../../config/db.yaml"
	cfg := config.MustLoadDBConfig(path)
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", cfg.User, cfg.Password, cfg.Port, cfg.DBName)
	ctx := context.Background()
	s := MustNewConnection(ctx, dbURL)

	id, err := s.SaveUser(ctx, "test1@gmail.com", []byte("qwertyy"))

	require.NoError(t, err)
	assert.Greater(t, int(id), 0)
}

// TestCannotSaveUserWithDuplicateEmail проверяет,
// что добавить юзера с уже существующим email в таблицу нельзя.
func TestCannotSaveUserWithDuplicateEmail(t *testing.T) {
	path := "../../../config/db.yaml"
	cfg := config.MustLoadDBConfig(path)
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", cfg.User, cfg.Password, cfg.Port, cfg.DBName)
	ctx := context.Background()
	s := MustNewConnection(ctx, dbURL)

	id, err := s.SaveUser(ctx, "TestCannotSaveUserWithDuplicateEmail@gmail.com", []byte("qwertyy"))
	require.NoError(t, err)
	assert.Greater(t, int(id), 0)

	id, err = s.SaveUser(ctx, "TestCannotSaveUserWithDuplicateEmail@gmail.com", []byte("qwertyy"))
	assert.ErrorIs(t, errors.Unwrap(err), storage.ErrUserExists)
	assert.Equal(t, 0, int(id))
}

// TestGetUser проверяет, что
// запрос на получения юзера успешно выполняется.
func TestGetUser(t *testing.T) {
	path := "../../../config/db.yaml"
	cfg := config.MustLoadDBConfig(path)
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", cfg.User, cfg.Password, cfg.Port, cfg.DBName)
	ctx := context.Background()
	s := MustNewConnection(ctx, dbURL)
	email := "TestGetUser@gmail.com"
	pass := []byte("qwe")

	_, err := s.SaveUser(ctx, email, pass)
	require.NoError(t, err)
	user, err := s.GetUser(ctx, email)

	require.NoError(t, err)
	assert.Equal(t, user.Email, email)
	assert.Equal(t, user.PaswordHash, pass)
	assert.Greater(t, int(user.Id), 0)
}

// TestCannotGetUserBecouseItsNotExists проверяет,
// что если юзера нет в базе данных, то возвращается пустая модель
// юзера и ошибка ErrUserNotFound.
func TestCannotGetUserBecouseItsNotExists(t *testing.T) {
	path := "../../../config/db.yaml"
	cfg := config.MustLoadDBConfig(path)
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", cfg.User, cfg.Password, cfg.Port, cfg.DBName)
	ctx := context.Background()
	s := MustNewConnection(ctx, dbURL)

	user, err := s.GetUser(ctx, "aboba")

	assert.ErrorIs(t, errors.Unwrap(err), storage.ErrUserNotFound)
	assert.Equal(t, user, models.User{})
}

// TestIsAdmin проверяет, что если юзерн не является
// админом, то метод IsAdmin возвращает false.
func TestIsAdmin(t *testing.T) {
	path := "../../../config/db.yaml"
	cfg := config.MustLoadDBConfig(path)
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", cfg.User, cfg.Password, cfg.Port, cfg.DBName)
	ctx := context.Background()
	s := MustNewConnection(ctx, dbURL)
	user_id, err := s.SaveUser(ctx, "TestIsAdmin@gmail.com", []byte("qwe"))
	require.NoError(t, err)

	isAdmin, err := s.IsAdmin(ctx, user_id)
	require.NoError(t, err)
	assert.Equal(t, isAdmin, false)
}

// TestCannotCheckAdminBecouseUserNotExists проверяет, что
// если попытаться проверить админа у несуществующего пользователя,
// то вернется ошибка ErrUserNotFound.
func TestCannotCheckAdminBecouseUserNotExists(t *testing.T) {
	path := "../../../config/db.yaml"
	cfg := config.MustLoadDBConfig(path)
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", cfg.User, cfg.Password, cfg.Port, cfg.DBName)
	ctx := context.Background()
	s := MustNewConnection(ctx, dbURL)

	isAdmin, err := s.IsAdmin(ctx, 2)

	assert.ErrorIs(t, errors.Unwrap(err), storage.ErrUserNotFound)
	assert.Equal(t, isAdmin, false)
}
