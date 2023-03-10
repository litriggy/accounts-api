package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"accounts/api/pkg/config"
	"accounts/api/pkg/middleware"
	"accounts/api/pkg/route"
	"accounts/api/platform/logger"
	"accounts/api/platform/memcached"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Serve() {
	appCfg := config.AppCfg()

	logger.SetUpLogger()
	logr := logger.GetLogger()

	app := fiber.New()
	app.Use(cors.New())
	middleware.FiberMiddleware(app)
	memcached.Init()

	route.SwaggerRoute(app)
	route.PublicRoutes(app)
	route.PrivateRoutes(app)

	route.UserRoutes(app.Group("/api/v1/user"))
	route.AuthRoutes(app.Group("/api/v1/auth"))
	route.TransactionRoutes(app.Group("/api/v1/tx"))
	route.InfoRouter(app.Group("/api/v1/info"))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sigCh
		logr.Infoln("Shutting down server")
		app.Shutdown()
	}()

	serverAddr := fmt.Sprintf("%s:%d", appCfg.Host, appCfg.Port)
	if err := app.Listen(serverAddr); err != nil {
		logr.Errorf("server is not running! error: %v", err)
	}
}
