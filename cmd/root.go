package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Глобальные переменные
var (
	kubeconfig string
	namespace  string
)

// initLogger инициализирует zerolog с учетом уровня логирования
func initLogger(levelStr string) {
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)

	var writer io.Writer
	if level <= zerolog.DebugLevel {
		writer = os.Stdout // JSON-лог
	} else {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}
	}

	log.Logger = zerolog.New(writer).
		With().
		Timestamp().
		Str("env", "dev").
		Str("version", "v0.1.0").
		Logger()
}

// rootCmd — корневая команда CLI
var rootCmd = &cobra.Command{
	Use:   "k8s-controller-tutorial",
	Short: "CLI with Cobra + Zerolog + Viper",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		level := viper.GetString("log-level")
		namespace = viper.GetString("namespace")
		initLogger(level)
	},

	Run: func(cmd *cobra.Command, args []string) {
		log.Trace().Msg("trace")
		log.Debug().Msg("debug")
		log.Info().Msg("info")
		log.Warn().Msg("warn")
		log.Error().Msg("error")
		fmt.Println("Welcome!")
	},
}

func init() {
	// === Viper настройки ===
	viper.SetConfigName("config") // config.yaml
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.k8s-controller-tutorial")
	viper.AutomaticEnv() // переменные окружения > config.yaml

	// Одинарный вывод конфиг-файла
	cobra.OnInitialize(func() {
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	})

	// Лог уровень
	rootCmd.PersistentFlags().String("log-level", "info", "set log level: trace, debug, info, warn, error")
	if err := viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level")); err != nil {
		log.Fatal().Err(err).Msg("failed to bind log-level flag")
	}

	// Kubeconfig путь
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file")

	// Namespace
	rootCmd.PersistentFlags().StringVar(&namespace, "namespace", "default", "Kubernetes namespace to use")
	if err := viper.BindPFlag("namespace", rootCmd.PersistentFlags().Lookup("namespace")); err != nil {
		log.Fatal().Err(err).Msg("failed to bind namespace flag")
	}
}

// Execute запускает CLI
func Execute() error {
	return rootCmd.Execute()
}
