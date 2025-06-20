// package cmd

// import (
// 	"fmt"

// 	"github.com/fasthttp/router"
// 	"github.com/rs/zerolog/log"
// 	"github.com/spf13/cobra"
// 	"github.com/spf13/viper"
// 	"github.com/valyala/fasthttp"
// )

// var port int

// var serverCmd = &cobra.Command{
// 	Use:   "server",
// 	Short: "Start a FastHTTP server",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		startFastHTTPServer()
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(serverCmd)
// 	serverCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
// 	viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
// }

// func startFastHTTPServer() {
// 	r := router.New()

// 	// простой обработчик
// 	r.GET("/", func(ctx *fasthttp.RequestCtx) {
// 		log.Info().Str("path", string(ctx.Path())).Msg("Handled request")
// 		fmt.Fprintf(ctx, "Hello from FastHTTP on port %d!\n", viper.GetInt("port"))
// 	})

//		addr := fmt.Sprintf(":%d", viper.GetInt("port"))
//		log.Info().Msgf("Starting FastHTTP server on %s", addr)
//		if err := fasthttp.ListenAndServe(addr, r.Handler); err != nil {
//			log.Fatal().Err(err).Msg("Server failed")
//		}
//	}
package cmd

import (
	"fmt"
	"time"

	"github.com/fasthttp/router"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

var port int

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a FastHTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		startFastHTTPServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
	viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
}

func startFastHTTPServer() {
	r := router.New()

	// Простой маршрут
	r.GET("/", logMiddleware(func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "Hello from FastHTTP on port %d!\n", viper.GetInt("port"))
	}))

	addr := fmt.Sprintf(":%d", viper.GetInt("port"))
	log.Info().Msgf("Starting FastHTTP server on %s", addr)
	if err := fasthttp.ListenAndServe(addr, r.Handler); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}

func logMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()

		// Request ID (если есть)
		requestID := string(ctx.Request.Header.Peek("X-Request-ID"))
		if requestID == "" {
			requestID = uuid.New().String()
			ctx.Response.Header.Set("X-Request-ID", requestID)
		}

		// Вызов обработчика
		next(ctx)

		duration := time.Since(start)

		log.Info().
			Str("method", string(ctx.Method())).
			Str("path", string(ctx.Path())).
			Str("remote_ip", ctx.RemoteIP().String()).
			Str("request_id", requestID).
			Dur("latency", duration).
			Msg("Request handled")
	}
}
