package main

import (
	"gRPC_sso/sso/internal/config"
	"gRPC_sso/sso/internal/lib/logger/handlers/slogpretty"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	// TODO: инициализировать объект конфига
	cfg := config.MustLoad()
	//fmt.Println("Конфиг инициализирован!\n", cfg)

	// TODO: инициализировать логгер (подключу slog из коробки go, можно подключить zap или logrus)
	log := setupLogger(envLocal)

	log.Info("env", slog.String("value", cfg.Env))
	log.Info("cfg", slog.Any("value", cfg))
	log.Info("port", slog.Int("value", cfg.GRPC.Port))

	log.Debug("debug message")
	log.Error("error message")
	log.Warn("warn message")

	// TODO: инициализировать приложение

	// TODO: запустить gPRC-сервер приложения

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	// Для локалки выводим в текст
	// Для прода выводим в JSON (ботам читать проще)
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

// Кастомный логгер для локальной разработки.
func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
