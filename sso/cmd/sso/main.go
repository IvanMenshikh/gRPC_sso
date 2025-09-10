package main

import (
	"gRPC_sso/sso/internal/app"
	"gRPC_sso/sso/internal/config"
	"gRPC_sso/sso/internal/lib/logger/handlers/slogpretty"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
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

	//log.Info("env", slog.String("value", cfg.Env))
	log.Info("cfg", slog.Any("value", cfg))
	//log.Info("port", slog.Int("value", cfg.GRPC.Port))

	//log.Debug("debug message")
	//log.Error("error message")
	//log.Warn("warn message")

	// TODO: инициализировать приложение
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	go application.GRPCSrv.MustRun()

	// Ждем сигналы завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Если сигнал пришел, то читаем и завершаем работу
	sgn := <-stop
	log.Info("stopping application", slog.String("signal", sgn.String()))
	
	application.GRPCSrv.Stop()
	log.Info("application stopped")
	// TODO: запустить gPRC-сервер приложения

}

// Устанавливаем логгер в зависимости от окружения
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
