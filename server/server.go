package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	access_control_api "github.com/yinloo-ola/tt-app/services/access_control/api"
	home "github.com/yinloo-ola/tt-app/services/home/api"
	"github.com/yinloo-ola/tt-app/util/template_util"
	"github.com/yinloo-ola/tt-app/views"
)

func main() {
	initLogger()

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("views/assets", false)))

	env := os.Getenv("GIN_MODE")

	var templateExecutor template_util.TemplateExecutor
	if env == "debug" {
		slog.Debug("DEV mode")
		router.LoadHTMLGlob(views.Glob())
		templateExecutor = &template_util.DebugTemplateExecutor{
			Glob: views.Glob(),
		}
	} else {
		gin.SetMode(gin.ReleaseMode)
		router.SetHTMLTemplate(views.ParseFS())
		templateExecutor = &template_util.ReleaseTemplateExecutor{
			Template: views.ParseGlob(),
		}
	}

	homeGroup := router.Group("/")
	home.AddAPIs(homeGroup, templateExecutor)

	accessControlGroup := router.Group("/access_control")
	access_control_api.AddAPIs(accessControlGroup, templateExecutor)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
		ErrorLog: slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		}), slog.LevelError),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
