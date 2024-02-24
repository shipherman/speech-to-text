package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"

	"github.com/shipherman/speech-to-text/gen/ent"
	"github.com/shipherman/speech-to-text/internal/jwt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Email struct {
	string
}
type Auth struct {
	// log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	tokenTTL    time.Duration
	secret      string
}

// Defined Errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Headers
var (
	headerAuthorize       = "authorization"
	headerEmail     Email = Email{"email"}
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
	jwtv5.RegisteredClaims

	// Username defines the identity of the user.
	Email string `json:"email"`
}

func New(
	// log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTL time.Duration,
	secret string,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		tokenTTL:    tokenTTL,
		secret:      secret,
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

	token, err := jwt.NewToken(*user, a.tokenTTL, a.secret)
	if err != nil {
		// a.log.Error("failed to generate token")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	s, _ := a.GetEmail(ctx)

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

// GetEmail returns email executed from JWT header
func (a *Auth) GetEmail(ctx context.Context) (string, error) {
	claims := &Claims{}

	tokenString, err := extractHeader(ctx, headerAuthorize)
	if err != nil {
		return "", err
	}

	token, err := jwtv5.ParseWithClaims(tokenString, claims, func(t *jwtv5.Token) (interface{}, error) {
		return []byte(a.secret), nil
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
func (a *Auth) AuthStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	ctx := stream.Context()
	fmt.Println("stream interceptor auth")
	err := a.CheckAuth(ctx)
	if err != nil {
		return err
	}

	err = handler(srv, stream)
	if err != nil {
		return err
	}
	return nil
}

// CheckAuth validates user token.
// Return error on invalid token
func (a *Auth) CheckAuth(ctx context.Context) (err error) {
	claims := &Claims{}
	tokenString, err := extractHeader(ctx, headerAuthorize)
	if err != nil {
		return err
	}
	//
	// !!! Hide secret key !!!
	token, err := jwtv5.ParseWithClaims(tokenString, claims, func(t *jwtv5.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
	if err != nil {
		log.Println(err)
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// extractHeader returns value for specified header
func extractHeader(ctx context.Context, header string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no headers in request")
	}

	authHeaders, ok := md[header]
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no header in request")
	}

	if len(authHeaders) != 1 {
		return "", status.Error(codes.Unauthenticated, "more than 1 header in request")
	}

	return authHeaders[0], nil
}
