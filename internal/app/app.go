package app

import (
	"fmt"
	"net/http"
	"time"

	"os"
	"os/signal"
	"social/internal/config"
	"social/internal/database/postgresql"
	"social/internal/service"
	"social/internal/transport/REST/auth"
	"social/internal/transport/REST/dialog"
	"social/internal/transport/REST/friend"
	"social/internal/transport/REST/middleware"
	"social/internal/transport/REST/post"
	"social/internal/transport/REST/user"
	"social/pkg/httpserver"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint: true,
	})
	log.SetReportCaller(true)

	db, err := postgresql.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	authSvc := service.NewAuthSvc(db)
	userSvc := service.NewUserSvc(db)
	authMiddleware := middleware.New(authSvc)

	r := gin.New()
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{}))

	r.POST("/login", auth.HandleLogin(authSvc))
	r.POST("/user/register", auth.HandleRegister(authSvc))

	r.Use(authMiddleware.Identify)

	r.GET("/user/get/:user_id", user.HandleGetUser(userSvc))
	r.GET("/user/search", user.HandleSearch())

	r.PUT("/friend/set/:user_id", friend.HandleSetFriend())
	r.PUT("/friend/delete/:user_id", friend.HandleDeleteFriend())

	r.POST("/post/create", post.HandleCreatePost())
	r.PUT("/post/update", post.HandleUpdatePost())
	r.PUT("/post/delete/:post_id", post.HandleDeletePost())
	r.GET("/post/get/:post_id", post.HandleGetPost())
	r.GET("/post/feed", post.HandleFeed())

	r.POST("/dialog/:user_id/send", dialog.HandleSend())
	r.GET("/dialog/:user_id/list", dialog.HandleList())

	runServer(r, cfg)
}

func runServer(handler http.Handler, cfg *config.Config) {
	readTimeout := time.Second * time.Duration(cfg.HTTPServer.ReadTimeout)
	writeTimeout := time.Second * time.Duration(cfg.HTTPServer.WriteTimeout)
	shutdownDuration := time.Second * time.Duration(cfg.HTTPServer.ShutdownTimeout)

	httpServer := httpserver.New(
		handler,
		httpserver.Addr(cfg.HTTPServer.Address()),
		httpserver.ReadTimeout(readTimeout),
		httpserver.WriteTimeout(writeTimeout),
		httpserver.ShutdownTimeout(shutdownDuration),
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	if err := httpServer.Shutdown(); err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	log.Info("server shutdown completed")
}
