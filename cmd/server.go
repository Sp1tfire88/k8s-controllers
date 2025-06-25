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
	// viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
	if err := viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level")); err != nil {
		log.Fatal().Err(err).Msg("failed to bind log-level flag")
	}
}

// func startFastHTTPServer() {
// 	r := router.New()

// 	// Обрабатывает GET и POST
// 	r.ANY("/", logMiddleware(handler))

// 	addr := fmt.Sprintf(":%d", viper.GetInt("port"))
// 	log.Info().Msgf("Starting FastHTTP server on %s", addr)
// 	if err := fasthttp.ListenAndServe(addr, r.Handler); err != nil {
// 		log.Fatal().Err(err).Msg("Server failed")
// 	}
// }

func startFastHTTPServer() {
	// 💡 Вызов информера ДО запуска HTTP-сервера
	go func() {
		err := StartDeploymentInformer(kubeconfig, inCluster, namespace)
		if err != nil {
			log.Fatal().Err(err).Msg("❌ Failed to start deployment informer")
		}
	}()

	// HTTP routes
	r := router.New()
	r.GET("/", logMiddleware(homeHandler))
	r.POST("/post", logMiddleware(postHandler))
	r.GET("/health", logMiddleware(healthHandler))

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

// func handler(ctx *fasthttp.RequestCtx) {
// 	switch string(ctx.Method()) {
// 	case fasthttp.MethodPost:
// 		body := ctx.PostBody()
// 		log.Info().
// 			Str("method", "POST").
// 			Str("path", string(ctx.Path())).
// 			Bytes("body", body).
// 			Msg("Received POST")

// 		ctx.SetStatusCode(fasthttp.StatusOK)
// 		ctx.SetBodyString("POST received")

// 	case fasthttp.MethodGet:
// 		log.Info().
// 			Str("method", "GET").
// 			Str("path", string(ctx.Path())).
// 			Msg("Handled GET")

// 		ctx.SetStatusCode(fasthttp.StatusOK)
// 		ctx.SetBodyString("Hello from FastHTTP")

// 	default:
// 		ctx.Error("Method Not Allowed", fasthttp.StatusMethodNotAllowed)
// 	}
// }

// обработчики маршрутов
func homeHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Welcome to the FastHTTP server!")
}

func postHandler(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	log.Info().Bytes("body", body).Msg("Received POST data")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("POST received")
}

func healthHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("OK")
}

// func userHandler(ctx *fasthttp.RequestCtx) {
// 	userID := ctx.UserValue("id").(string)
// 	ctx.SetStatusCode(fasthttp.StatusOK)
// 	ctx.SetBodyString(fmt.Sprintf("User ID: %s", userID))
// }
