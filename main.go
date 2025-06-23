package main

import (
	"github.com/Sp1tfire88/k8s-controllers/cmd"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("command execution failed")
	}
}
