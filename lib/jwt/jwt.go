package jwt

import (
	"sso/interanal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, app models.App, ttl time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(ttl).Unix()
	claims["app_id"] = app.Id

	signedToken, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
