package Tool

import (
	"context"
	"errors"
	"net/http"
	"reflect"

	"github.com/golang-jwt/jwt"
)

var ErrInvalidToken = errors.New("invalid token")

//l.svcCtx.Config.Auth.
func UserIdFromCookie(r *http.Request, AccessSecret string) (string, error) {
	cookie, err := r.Cookie("Authorization")
	if err != nil || cookie == nil || cookie.Value == "" {
		return "", err
	}
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(AccessSecret), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidToken
	}

	idInterface := claims["id"]
	if idInterface == nil || reflect.TypeOf(idInterface).Kind() != reflect.String {
		return "", ErrInvalidToken
	}
	if value, ok := idInterface.(string); ok {
		return value, nil
	}
	return "", ErrInvalidToken
}

var ZeroUID error = errors.New("zero uid")

func UserIdFromContext(ctx context.Context) (uid string, err error) {
	var (
		uidInterface interface{}
		ok           bool
	)
	if uidInterface = ctx.Value("id"); uidInterface == nil {
		return "", ZeroUID
	}

	if uid, ok = ctx.Value("id").(string); ok == false {
		return "", ZeroUID
	}
	return uid, nil
}

//validate user id
func ValidateJwtUser(ctx context.Context, uid string) bool {
	UId, err := UserIdFromContext(ctx)
	if err != nil {
		return false
	}
	return UId == uid
}
