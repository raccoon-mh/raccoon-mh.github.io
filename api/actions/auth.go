package actions

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"api/util"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/golang-jwt/jwt"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	// Here you can handle the user info, e.g., save it to the database
	// tx := c.Value("tx").(*pop.Connection)
	// You would create or find the user in the database here

	fmt.Println(user)

	// JWT 생성
	token, err := generateJWT(user)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	commonRes := util.CommonResponseStatusOK(token)
	return c.Render(commonRes.Status.StatusCode, r.JSON(commonRes))
}

func generateJWT(user goth.User) (string, error) {
	secret := envy.Get("JWT_SECRET", "default_secret_key")
	claims := jwt.MapClaims{
		"id":       user.UserID,
		"name":     user.Name,
		"email":    user.Email,
		"provider": user.Provider,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // 72시간 유효
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func parseJWT(tokenStr string) (jwt.MapClaims, error) {
	secret := envy.Get("SESSION_SECRET", "raccoon-mh")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func AuthMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		cookie, err := c.Request().Cookie("auth_token")
		if err != nil {
			return c.Error(http.StatusUnauthorized, errors.New("unauthorized"))
		}
		claims, err := parseJWT(cookie.Value)
		if err != nil {
			return c.Error(http.StatusUnauthorized, errors.New("unauthorized"))
		}
		c.Set("user", claims)
		return next(c)
	}
}
