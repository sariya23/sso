package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sso/interanal/domain/models"
	"sso/interanal/storage"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Storage struct {
	connection *pgx.Conn
}

func MustNewConnection(ctx context.Context, dbURL string) *Storage {
	const op = "storage.postgres.MustNewConnection"
	ctx, cancel := context.WithTimeout(ctx, time.Second*4)
	defer cancel()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("%s: cannot connect to db with URL: %s, with error: %v", op, dbURL, err)
	}
	err = conn.Ping(ctx)
	if err != nil {
		log.Fatalf("%s: db is unreachable: %v", op, err)
	}
	return &Storage{connection: conn}
}

func (s *Storage) Stop(ctx context.Context) {
	const op = "storage.postgres.Stop"
	err := s.connection.Close(ctx)
	if err != nil {
		log.Fatalf("%s: cannot close db connection: %v", op, err)
	}
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"
	var pgErr *pgconn.PgError
	var userId int64
	stmt := `insert into "user"(email, passHash) values ($1, $2) returning user_id`
	err := s.connection.QueryRow(ctx, stmt, email, passHash).Scan(&userId)
	if err != nil {
		if ok := errors.As(err, &pgErr); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}

func (s *Storage) GetUser(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.GetUser"
	type Row struct {
		id       int
		email    string
		passHash []byte
	}
	var r Row
	stmt := `select user_id, email, pass_hash from "user" where email=$1`
	err := s.connection.QueryRow(ctx, stmt, email).Scan(&r.id, r.email, r.passHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, storage.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.User{Id: int64(r.id), Email: r.email, PaswordHash: r.passHash}, nil
}

func (s *Storage) GetApp(ctx context.Context, appId int) (models.App, error) {
	const op = "storage.postgres.GetApp"
	type Row struct {
		id     int
		name   string
		secret string
	}
	var r Row
	stmt := `select app_id, name, secret from app where app_id=$1`
	err := s.connection.QueryRow(ctx, stmt, appId).Scan(r.id, r.name, r.secret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return models.App{Id: r.id, Name: r.name, Secret: r.secret}, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"
	var isAdmin bool
	stmt := `select id_admin from "users" where user_id=$1`
	err := s.connection.QueryRow(ctx, stmt, userId).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return isAdmin, nil
}
