// package main

// import (
// 	"os"

// 	"github.com/Sp1tfire88/k8s-controllers/cmd"
// 	"github.com/Sp1tfire88/k8s-controllers/pkg/logger"
// )

// func main() {
// 	// Проверь аргументы (например, --debug)
// 	debug := false
// 	for _, arg := range os.Args {
// 		if arg == "--debug" {
// 			debug = true
// 			break
// 		}
// 	}

// 	// Инициализируй логгер
// 	logger.Setup(debug)

// 	// Запуск CLI-команды
// 	cmd.Execute()
// }

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import "github.com/Sp1tfire88/k8s-controllers/cmd"

func main() {
	cmd.Execute()
}
