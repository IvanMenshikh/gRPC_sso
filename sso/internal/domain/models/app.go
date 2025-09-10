package models

// Модель приложения
// secret - для подписи JWT токенов
type App struct {
	ID  int
	Name string
	Secret string
}

