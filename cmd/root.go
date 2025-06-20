package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Глобальная переменная для флага debug
var debug bool

// Функция инициализации логгера с ConsoleWriter
func initLogger(debug bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

// Корневая команда
var rootCmd = &cobra.Command{
	Use:   "k8s-controller-tutorial",
	Short: "A brief description of your application",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initLogger(debug) // инициализируем логгер перед выполнением любой команды
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("This is an info log")
		log.Debug().Msg("This is a debug log")
		log.Trace().Msg("This is a trace log")
		log.Warn().Msg("This is a warn log")
		log.Error().Msg("This is an error log")

		fmt.Println("Welcome to k8s-controller-tutorial CLI!")
	},
}

func init() {
	// Регистрируем глобальный флаг --debug
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")
}

func Execute() error {
	return rootCmd.Execute()
}
