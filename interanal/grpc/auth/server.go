package auth

import (
	"context"
	"errors"
	"net/mail"
	"sso/interanal/service/auth"
	"sso/interanal/storage"

	ssov1 "github.com/sariya23/sso_proto/gen/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyAppId = 0
)

type userCreds struct {
	email, password string
}

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appId int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userId int64, err error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type ServerAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func RegisterServerAPI(grpcServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(grpcServer, &ServerAPI{auth: auth})
}

func (s *ServerAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()
	if err := validateUserCreds(userCreds{email: email, password: password}); err != nil {
		return nil, err
	}
	if req.GetAppId() == emptyAppId {
		return nil, status.Error(codes.InvalidArgument, "app id is required")
	}

	token, err := s.auth.Login(ctx, email, req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCreds) {
			return nil, status.Error(codes.InvalidArgument, "invalid creds")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	email, password := req.GetEmail(), req.GetPassword()
	if err := validateUserCreds(userCreds{email: email, password: password}); err != nil {
		return nil, err
	}
	userId, err := s.auth.RegisterNewUser(ctx, email, password)
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, auth.ErrUserExists.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.RegisterResponse{
		UserId: userId,
	}, nil
}

func (s *ServerAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	userId := req.GetUserId()
	if userId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}
	isAdmin, err := s.auth.IsAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateUserCreds(creds userCreds) error {
	if _, err := mail.ParseAddress(creds.email); err != nil {
		return status.Error(codes.InvalidArgument, "email is invalid")
	}
	if creds.password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	return nil
}
