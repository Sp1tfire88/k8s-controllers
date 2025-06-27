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
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/Sp1tfire88/k8s-controllers/pkg/controller"
)

var (
	port            int
	leaderElection  bool
	metricsPort     int
	configNamespace string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a FastHTTP server with controller-runtime",
	Run: func(cmd *cobra.Command, args []string) {
		startFastHTTPServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// CLI flags (все можно задать через конфиг)
	serverCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
	serverCmd.Flags().BoolVar(&leaderElection, "enable-leader-election", true, "Enable leader election for controller manager")
	serverCmd.Flags().IntVar(&metricsPort, "metrics-port", 8081, "Port for controller manager metrics endpoint")
	serverCmd.Flags().StringVar(&configNamespace, "namespace", "default", "Kubernetes namespace to watch")

	// Bind к viper (поддержка конфига и CLI одновременно)
	_ = viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
	_ = viper.BindPFlag("enableLeaderElection", serverCmd.Flags().Lookup("enable-leader-election"))
	_ = viper.BindPFlag("metricsPort", serverCmd.Flags().Lookup("metrics-port"))
	_ = viper.BindPFlag("namespace", serverCmd.Flags().Lookup("namespace"))
	_ = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
}

func startFastHTTPServer() {
	// --- 0. Загрузить config.yaml, если есть ---
	if configFile := viper.ConfigFileUsed(); configFile == "" {
		// По умолчанию читаем config.yaml из cwd
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		_ = viper.ReadInConfig() // не fail, если нет файла
	}

	// --- 1. Считываем параметры ---
	metricsPort := viper.GetInt("metricsPort")
	leaderElection := viper.GetBool("enableLeaderElection")
	namespace := viper.GetString("namespace")

	// --- 2. Setup controller-runtime logger ---
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// --- 3. Setup controller-runtime manager ---
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(appsv1.AddToScheme(scheme))

	cfg := ctrl.GetConfigOrDie()
	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme,
		Metrics: server.Options{
			BindAddress: fmt.Sprintf(":%d", metricsPort),
		},
		LeaderElection:   leaderElection,
		LeaderElectionID: "k8s-controllers-lock",
		// Namespace:     // НЕ поддерживается в v0.18.2!
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start controller-runtime manager")
	}

	// --- 4. Register the Deployment controller с фильтрацией по namespace ---
	if err := (&controller.DeploymentReconciler{
		Client:    mgr.GetClient(),
		Scheme:    mgr.GetScheme(),
		Namespace: namespace, // добавим это поле в структуру
	}).SetupWithManager(mgr); err != nil {
		log.Fatal().Err(err).Msg("Failed to setup Deployment controller")
	}

	// --- 5. Start manager в фоне ---
	go func() {
		log.Info().Msg("🔧 Starting controller-runtime manager")
		if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
			log.Fatal().Err(err).Msg("controller-runtime manager exited")
		}
	}()

	// --- 6. FastHTTP server ---
	r := router.New()
	r.GET("/", logMiddleware(homeHandler))
	r.POST("/post", logMiddleware(postHandler))
	r.GET("/health", logMiddleware(healthHandler))
	r.GET("/deployments", logMiddleware(deploymentsHandler))

	addr := fmt.Sprintf(":%d", viper.GetInt("port"))
	log.Info().Msgf("🚀 Starting FastHTTP server on %s", addr)
	if err := fasthttp.ListenAndServe(addr, r.Handler); err != nil {
		log.Fatal().Err(err).Msg("FastHTTP server failed")
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
