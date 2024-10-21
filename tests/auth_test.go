package tests

import (
	"errors"
	"sso/interanal/service/auth"
	"sso/tests/suite"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	ssov1 "github.com/sariya23/sso_proto/gen/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppId     = 0
	appId          = 1
	appSecret      = "test-secret"
	passDefaultLen = 10
)

// TestSuccessLogin проверяет, что
//
// - если пользователь существует и он указал верные креды, то
// ему возвращается валидный jwt.
func TestSuccessLogin(t *testing.T) {
	ctx, suite := suite.New(t)
	email := gofakeit.Email()
	password := randomFakePssword()
	resRegister, err := suite.AuthClient.Register(
		ctx,
		&ssov1.RegisterRequest{
			Email:    email,
			Password: password,
		})
	require.NoError(t, err)
	assert.NotEmpty(t, resRegister.GetUserId())

	resLogin, err := suite.AuthClient.Login(ctx, &ssov1.LoginRequest{Email: email, Password: password, AppId: appId})
	require.NoError(t, err)
	loginTime := time.Now()
	token := resLogin.GetToken()
	assert.NotEmpty(t, token)
	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)
	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, resRegister.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appId, int(claims["app_id"].(float64)))
	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(suite.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

// TestUserCannotRegiterTwice проверяет, что
// пользователь с одним и тем же email
// не может зарегистрироваться дважды.
func TestUserCannotRegiterTwice(t *testing.T) {
	ctx, suite := suite.New(t)
	email := gofakeit.Email()
	password := randomFakePssword()
	resp, err := suite.AuthClient.Register(
		ctx,
		&ssov1.RegisterRequest{
			Email:    email,
			Password: password,
		},
	)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.GetUserId())

	resp, err = suite.AuthClient.Register(
		ctx,
		&ssov1.RegisterRequest{
			Email:    email,
			Password: password,
		},
	)
	require.ErrorIs(t, errors.Unwrap(err), auth.ErrUserExists)
	require.Equal(t, resp.GetUserId(), int64(0))
}

func randomFakePssword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
