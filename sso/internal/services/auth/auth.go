package auth

import (
	"context"
	"gRPC_sso/sso/internal/domain/models"
	"log/slog"
	"time"
)

// Сервис аутентификации
type Auth struct {
	log          *slog.Logger
	usersaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

// Для гибкости разделим интерфейсы
// Один сохраняет пользователя (UserSaver)
// Второй предоставляет данные о пользователе (UserProvider)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

// Нужен для получения секрета приложения по его ID
type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

// New создает новый экземпляр Auth с заданными зависимостями.
//
// Принимает:
//   - log: логгер для записи событий и ошибок,
//   - userSaver: интерфейс для сохранения пользователей,
//   - userProvider: интерфейс для получения данных пользователей,
//   - appProvider: интерфейс для работы с приложениями,
//   - tokenTTL: время жизни токена JWT для авторизации.
//
// Возвращает указатель на Auth, готовый к использованию.
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log: log, 
		usersaver: userSaver,
		userProvider: userProvider,
		appProvider: appProvider,
		tokenTTL: tokenTTL,
	}
}

// Login выполняет аутентификацию пользователя по email и паролю.
//
// Принимает:
//   - ctx: контекст запроса,
//   - email: email пользователя,
//   - password: пароль пользователя,
//   - appID: идентификатор приложения или клиента, из которого пришёл запрос.
//
// Возвращает:
//   - строку (например, токен доступа, JWT или session ID),
//   - error
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, error){
	panic("implement me")
}

// RegisterNewUser создает нового пользователя в системе.
// Выполняет валидацию данных и сохраняет запись в хранилище.
//
// Принимает:
//   - ctx: контекст запроса,
//   - email: адрес электронной почты пользователя,
//   - password: пароль пользователя (в открытом виде).
//
// Возвращает:
//   - userID: идентификатор созданного пользователя,
//   - error
// В случае ошибки валидации или сбоя при сохранении возвращает error.
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	panic("implement me")
}

// IsAdmin проверяет, является ли пользователь с заданным идентификатором администратором.
//
// Принимает:
//   - ctx: контекст запроса (для отмены или таймаутов),
//   - userID: идентификатор пользователя, которого проверяем.
//
// Возвращает:
//   - true, если пользователь имеет права администратора, иначе false,
//   - error, если произошла ошибка при проверке (например, проблема с хранилищем данных).
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	panic("implement me")
}


