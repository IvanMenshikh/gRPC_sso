package auth

import (
	"context"

	ssov1 "github.com/IvanMenshikh/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValue = 0
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

// Регистрируем обработчик
func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{
		auth: auth,
	})
}

// gRPC ручки

// Register создает нового пользователя.
//
// Проверяет входные данные через validateRegister и вызывает сервис Auth для создания пользователя.
//
// Возвращает идентификатор нового пользователя или error.
func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		// TODO: Дописать обработку ошибок
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

// Login выполняет аутентификацию пользователя по email и паролю.
//
// Проверяет входные данные через validateLogin и вызывает сервис Auth для получения токена.
//
// Возвращает JWT/токен или error.
func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	// TODO: Имплементируем логин через сервис авторизации
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		// TODO: Дописать обработку ошибок
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

// IsAdmin проверяет, является ли пользователь администратором.
//
// Проверяет входные данные через validateIsAdmin и вызывает сервис Auth.
//
// Возвращает true/false или error.
func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		// TODO: Дописать обработку ошибок
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil

}

// Валидация данных вручную!

// validateLogin проверяет корректность данных для запроса Login.
//
// Проверяет:
//   - email не пустой,
//   - password не пустой и минимум 12 символов,
//   - app_id установлен (не равен emptyValue).
//
// В случае ошибки возвращает gRPC статус codes.InvalidArgument.
func validateLogin(req *ssov1.LoginRequest) error {
	// Проверяем на пустые или нулевые значения
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	// Проверяем на пустые или нулевые значения
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	// Проверяем на минимальную длину пароля
	if len(req.GetPassword()) < 12 {
		return status.Error(codes.InvalidArgument, "password must be at least 12 characters long")
	}
	// Проверяем на пустые или нулевые значения
	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}
	return nil
}

// validateRegister проверяет корректность данных для запроса Register.
//
// Проверяет, что email и password не пустые.
//
// В случае ошибки возвращает gRPC статус codes.InvalidArgument.
func validateRegister(req *ssov1.RegisterRequest) error {

	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

// validateIsAdmin проверяет корректность данных для запроса IsAdmin.
//
// Проверяет, что user_id установлен (не равен emptyValue).
//
// В случае ошибки возвращает gRPC статус codes.InvalidArgument.
func validateIsAdmin(req *ssov1.IsAdminRequest) error {

	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}
