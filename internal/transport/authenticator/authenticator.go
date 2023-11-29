package authenticator

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/shipherman/speech-to-text/gen/ent"
	"github.com/shipherman/speech-to-text/gen/ent/user"
	"github.com/shipherman/speech-to-text/internal/models"
)

type Claims struct {
	jwt.RegisteredClaims
	User string
}

type Authenticator struct {
	Client ent.Client
}

const tockenExpiration = time.Hour * 3
const sercretKey = "supersecretkey"

var ErrUserDoesNotExist = errors.New("no such user")

func NewAuthenticator(dbclient ent.Client) Authenticator {
	return Authenticator{Client: dbclient}
}

func buildJWTString(user string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tockenExpiration)),
		},
		User: user,
	})
	tokenString, err := token.SignedString([]byte(sercretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getUser(tokenString string) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(sercretKey), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return claims.User, fmt.Errorf("invalid token")
	}

	return claims.User, nil
}

// Auth returns JWT string
func (a *Authenticator) Auth(u, p string) (jwt string, err error) {
	ctx := context.Background()

	exist, err := a.Client.User.Query().
		Where(user.LoginEQ(u)).
		Exist(ctx)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", ErrUserDoesNotExist
	}

	entUser, _ := a.Client.User.Query().
		Where(user.LoginEQ(u)).
		First(ctx)

	if entUser.Password != p {
		return "", fmt.Errorf("wrong password")
	}

	return buildJWTString(u)
}

func (a *Authenticator) CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if authenticated
		// Return 401 if not
		JWT := r.Header.Get("Authorization")
		if JWT == "" {
			http.Error(w, "AccessDenied", http.StatusUnauthorized)
			return
		}

		// Get user
		userJWT, err := getUser(JWT)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Add user as context parameter
		r = r.WithContext(context.WithValue(r.Context(), models.UserCtxKey{}, userJWT))

		next.ServeHTTP(w, r)
	})

}
