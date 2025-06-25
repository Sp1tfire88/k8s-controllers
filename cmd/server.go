package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fasthttp/router"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	appsv1 "k8s.io/api/apps/v1"
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
	if err := viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level")); err != nil {
		log.Fatal().Err(err).Msg("failed to bind log-level flag")
	}
}

func startFastHTTPServer() {
	if err := StartDeploymentInformerFromConfig(); err != nil {
		log.Warn().Err(err).Msg("Informer not started")
	}

	r := router.New()
	r.GET("/", logMiddleware(homeHandler))
	r.POST("/post", logMiddleware(postHandler))
	r.GET("/health", logMiddleware(healthHandler))
	r.GET("/deployments", logMiddleware(deploymentsHandler)) // ✅ Новый endpoint

	addr := fmt.Sprintf(":%d", viper.GetInt("port"))
	log.Info().Msgf("Starting FastHTTP server on %s", addr)
	if err := fasthttp.ListenAndServe(addr, r.Handler); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}

func logMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		requestID := string(ctx.Request.Header.Peek("X-Request-ID"))
		if requestID == "" {
			requestID = uuid.New().String()
			ctx.Response.Header.Set("X-Request-ID", requestID)
		}
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

// ✅ Новый обработчик JSON API
func deploymentsHandler(ctx *fasthttp.RequestCtx) {
	store := GetDeploymentStore()
	if store == nil {
		ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
		ctx.SetBodyString(`{"error":"deployment cache not ready"}`)
		return
	}

	var names []string
	for _, obj := range store.List() {
		if d, ok := obj.(*appsv1.Deployment); ok {
			names = append(names, d.GetName())
		}
	}

	resp, err := json.Marshal(names)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(`{"error":"failed to serialize deployments"}`)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(resp)
}
