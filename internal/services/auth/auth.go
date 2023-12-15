package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	dgjwt "github.com/dgrijalva/jwt-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/shipherman/speech-to-text/gen/ent"
	"github.com/shipherman/speech-to-text/internal/jwt"

	"google.golang.org/grpc"
)

type Auth struct {
	// log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	tokenTTL    time.Duration
	Secret      string
}

// Defined Errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Headers
var (
	headerAuthorize = "authorization"
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		username string,
		email string,
		pass string,
	) (uid int64, err error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (*ent.User, error)
}

// Claims defines the struct containing the token claims.
type Claims struct {
	dgjwt.StandardClaims

	// Username defines the identity of the user.
	Email string `json:"email"`
}

func New(
	// log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		tokenTTL:    tokenTTL,
	}
}

// Login checks if user with given credentials exists in the system and returns access token.
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error.
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	const op = "Auth.Login"

	// log := a.log.With(
	// 	slog.String("op", op),
	// 	slog.String("username", email),
	// )

	// log.Info("attempting to login user")

	user, err := a.usrProvider.GetUser(ctx, email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	if user.Password != password {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// log.Info("user logged in successfully")

	token, err := jwt.NewToken(*user, a.tokenTTL, a.Secret)
	if err != nil {
		// a.log.Error("failed to generate token")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	s, _ := a.GetEmail(context.Background(), token)

	fmt.Println(s)

	return token, nil
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (a *Auth) RegisterNewUser(ctx context.Context,
	username string,
	email string,
	pass string,
) (int64, error) {
	const op = "Auth.RegisterNewUser"

	// log := a.log.With(
	// 	slog.String("op", op),
	// 	slog.String("email", email),
	// )

	// log.Info("registering user")

	id, err := a.usrSaver.SaveUser(ctx, username, email, pass)
	if err != nil {
		// log.Error("failed to save user")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) GetEmail(ctx context.Context,
	tokenString string,
) (string, error) {
	claims := &Claims{}

	token, err := dgjwt.ParseWithClaims(tokenString, claims, func(t *dgjwt.Token) (interface{}, error) {
		return []byte(a.Secret), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return claims.Email, fmt.Errorf("invalid token")
	}

	return claims.Email, nil
}

// AuthInterceptor provides auth for api
func AuthUnaryInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	fmt.Println("checking auth")
	email, err := CheckAuth(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("gRPC method: %s, %v", info.FullMethod, req)

	newCtx := ctx

	if len(email) > 0 {
		newCtx = context.WithValue(ctx, "email", email)
		log.Println(newCtx.Value("email"))
	}

	return handler(newCtx, req)
}

// CheckAuth validates user token.
// Return error on invalid token
func CheckAuth(ctx context.Context) (email string, err error) {
	tokenStr := getTokenFromContext(ctx)
	if len(tokenStr) == 0 {
		return "", fmt.Errorf("empty user token")
	}

	// Define claims to parse user token
	var clientClaims Claims

	token, err := dgjwt.ParseWithClaims(tokenStr, &clientClaims, func(token *dgjwt.Token) (interface{}, error) {
		if token.Header["alg"] != "HS256" {
			return nil, fmt.Errorf("ErrInvalidAlgorithm")
		}
		return token, nil
	})
	if err != nil {
		return "", fmt.Errorf("jwt parse error")
	}

	if !token.Valid {
		fmt.Println("ErrInvalidToken")
	}

	return clientClaims.Email, nil
}

func getTokenFromContext(ctx context.Context) string {
	val := metautils.ExtractIncoming(ctx).Get(headerAuthorize)
	return val
}
