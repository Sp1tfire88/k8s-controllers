package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Глобальная переменная для флага уровня логирования
var logLevel string

// initLogger инициализирует zerolog с заданным уровнем логирования и форматированием
// func initLogger(levelStr string) {
// 	zerolog.TimeFieldFormat = time.RFC3339

// 	level, err := zerolog.ParseLevel(levelStr)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Invalid log level '%s', fallback to 'info'\n", levelStr)
// 		level = zerolog.InfoLevel
// 	}

// 	zerolog.SetGlobalLevel(level)

// 	log.Logger = zerolog.New(zerolog.ConsoleWriter{
// 		Out:        os.Stdout,
// 		TimeFormat: "15:04:05",
// 	}).With().
// 		Timestamp().
// 		Str("env", "dev").
// 		Str("version", "v0.1.0").
// 		Logger()
// }

func initLogger(levelStr string) {
	zerolog.TimeFieldFormat = time.RFC3339

	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid log level '%s', fallback to 'info'\n", levelStr)
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Файл для логов
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open log file: %v\n", err)
		logFile = os.Stdout
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
	}

	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)

	log.Logger = zerolog.New(multi).
		With().
		Timestamp().
		Str("env", "dev").
		Str("version", "v0.1.0").
		Logger()
}

// rootCmd — корневая команда CLI
var rootCmd = &cobra.Command{
	Use:   "k8s-controller-tutorial",
	Short: "A brief description of your application",
	Long:  "This is a CLI application for demonstrating Cobra and Zerolog integration.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// initLogger(logLevel)
		initLogger(viper.GetString("log-level"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Trace().Msg("This is a trace log")
		log.Debug().Msg("This is a debug log")
		log.Info().Msg("This is an info log")
		log.Warn().Msg("This is a warn log")
		log.Error().Msg("This is an error log")

		fmt.Println("Welcome to k8s-controller-tutorial CLI!")
	},
}

// init регистрирует глобальные флаги
// func init() {
// 	viper.AutomaticEnv() // автоматически считывает переменные окружения

// 	// Привязываем ENV к флагу
// 	viper.BindEnv("log-level", "LOG_LEVEL")

// 	// Привязываем флаг cobra -> viper
// 	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "set log level: trace, debug, info, warn, error")
// 	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))

// }

func init() {
	// === Viper настройки ===
	viper.SetConfigName("config")                         // имя файла без расширения
	viper.SetConfigType("yaml")                           // тип конфигурационного файла
	viper.AddConfigPath(".")                              // искать в текущей директории
	viper.AddConfigPath("$HOME/.k8s-controller-tutorial") // или ~/.k8s-controller-tutorial
	viper.AutomaticEnv()                                  // переменные окружения переопределяют

	// Попробуем прочитать файл конфигурации
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("No config file found or failed to read:", err)
	}

	// Привязка флагов
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "set log level: trace, debug, info, warn, error")
	viper.BindEnv("log-level", "LOG_LEVEL")
	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
}

// Execute запускает CLI
func Execute() error {
	return rootCmd.Execute()
}
