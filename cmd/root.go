// package cmd

// import (
// 	"fmt"
// 	"os"

// 	"github.com/rs/zerolog"
// 	"github.com/rs/zerolog/log"
// 	"github.com/spf13/cobra"
// )

// // Глобальная переменная для флага debug
// var debug bool

// // Функция инициализации логгера с ConsoleWriter
// func initLogger(debug bool) {
// 	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

// 	if debug {
// 		zerolog.SetGlobalLevel(zerolog.DebugLevel)
// 	} else {
// 		zerolog.SetGlobalLevel(zerolog.InfoLevel)
// 	}

// 	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
// }

// // Корневая команда
// var rootCmd = &cobra.Command{
// 	Use:   "k8s-controller-tutorial",
// 	Short: "A brief description of your application",
// 	PersistentPreRun: func(cmd *cobra.Command, args []string) {
// 		initLogger(debug) // инициализируем логгер перед выполнением любой команды
// 	},
// 	Run: func(cmd *cobra.Command, args []string) {
// 		log.Info().Msg("This is an info log")
// 		log.Debug().Msg("This is a debug log")
// 		log.Trace().Msg("This is a trace log")
// 		log.Warn().Msg("This is a warn log")
// 		log.Error().Msg("This is an error log")

// 		fmt.Println("Welcome to k8s-controller-tutorial CLI!")
// 	},
// }

// func init() {
// 	// Регистрируем глобальный флаг --debug
// 	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")
// }

// func Execute() error {
// 	return rootCmd.Execute()
// }

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Глобальная переменная для флага уровня логирования
var logLevel string

// initLogger инициализирует zerolog с заданным уровнем логирования и форматированием
func initLogger(levelStr string) {
	zerolog.TimeFieldFormat = time.RFC3339

	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid log level '%s', fallback to 'info'\n", levelStr)
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)

	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
	}).With().
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
		initLogger(logLevel)
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
func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "set log level: trace, debug, info, warn, error")
}

// Execute запускает CLI
func Execute() error {
	return rootCmd.Execute()
}
