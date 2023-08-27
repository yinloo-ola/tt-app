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
	accesscontrol "github.com/yinloo-ola/tt-app/services/access-control/api"
	home "github.com/yinloo-ola/tt-app/services/home/api"
	"github.com/yinloo-ola/tt-app/templates"
)

func main() {
	initLogger()

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("../templates/assets", false)))
	router.SetHTMLTemplate(templates.Parse())

	homeGroup := router.Group("/")
	home.AddAPIs(homeGroup)

	accessControlGroup := router.Group("/access_control")
	accesscontrol.AddAPIs(accessControlGroup)

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
