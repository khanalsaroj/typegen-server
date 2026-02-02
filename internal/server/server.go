package server

import (
	"context"
	"net/http"
	"time"

	"github.com/khanalsaroj/typegen-server/internal/modules/connection"
	"github.com/khanalsaroj/typegen-server/internal/modules/health"
	"github.com/khanalsaroj/typegen-server/internal/pkg/crypto"

	"github.com/khanalsaroj/typegen-server/internal/config"
	"github.com/khanalsaroj/typegen-server/internal/middleware"

	dbHandlerPkg "github.com/khanalsaroj/typegen-server/internal/modules/conn/handler"
	dbRepoPkg "github.com/khanalsaroj/typegen-server/internal/modules/conn/repository"
	dbServicePkg "github.com/khanalsaroj/typegen-server/internal/modules/conn/service"

	typeHandlerPkg "github.com/khanalsaroj/typegen-server/internal/modules/gentype/handler"
	typeServicePkg "github.com/khanalsaroj/typegen-server/internal/modules/gentype/service"

	mprHandlerPkg "github.com/khanalsaroj/typegen-server/internal/modules/mapper/handler"
	mprServicePkg "github.com/khanalsaroj/typegen-server/internal/modules/mapper/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	router    *gin.Engine
	server    *http.Server
	config    *config.Config
	db        *gorm.DB
	cryptoSvc *crypto.Service
	logger    *zap.Logger
}

func New(cfg *config.Config, db *gorm.DB, cryptoSvc *crypto.Service, logger *zap.Logger) *Server {
	if cfg.App.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	s := &Server{
		router:    router,
		config:    cfg,
		db:        db,
		cryptoSvc: cryptoSvc,
		logger:    logger,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

func (s *Server) setupMiddleware() {
	s.router.Use(middleware.Recovery(s.logger))
	s.router.Use(middleware.Logger(s.logger))
	s.router.Use(middleware.CORS(s.config.Security.CORSAllowOrigins))
	s.router.Use(middleware.SecureHeaders())

	if s.config.Security.RateLimitEnabled {
		s.router.Use(middleware.RateLimit(s.config.Security.RateLimitRPS))
	}
}

func (s *Server) setupRoutes() {

	v1 := s.router.Group("/api/v1")
	{
		healthService := health.NewService(s.db, s.config.App.Version)
		healthHandler := health.NewHandler(healthService)

		v1.GET("/health", healthHandler.Health)

		repo := dbRepoPkg.New()
		dbSvc := dbServicePkg.New(repo)
		handler := dbHandlerPkg.NewHandler(dbSvc)

		dbGroup := v1.Group("connection")
		{
			dbGroup.POST("/test", handler.Connect)
		}

		dbRepo := connection.NewRepository(s.db)
		dbService := connection.NewService(dbRepo, s.cryptoSvc)
		typeSvc := &typeServicePkg.TypeService{
			ConnectionService: dbService,
		}
		typeHandler := typeHandlerPkg.New(typeSvc)

		typeGroup := v1.Group("/type")
		{
			typeGroup.POST("", typeHandler.GenerateType)
		}

		mprSvc := &mprServicePkg.MprService{
			ConnectionService: dbService,
		}
		mprHandler := mprHandlerPkg.New(mprSvc)

		mprGroup := v1.Group("/mapper")
		{
			mprGroup.POST("", mprHandler.GenerateMapper)
		}

		userRepo := connection.NewRepository(s.db)
		userService := connection.NewService(userRepo, s.cryptoSvc)
		userHandler := connection.NewHandler(userService)

		connectionGroup := v1.Group("/connection")

		connectionGroup.POST("", userHandler.Create)
		connectionGroup.GET("/:id", userHandler.GetByID)
		connectionGroup.GET("", userHandler.List)
		connectionGroup.PUT("/:id", userHandler.Update)
		connectionGroup.DELETE("/:id", userHandler.Delete)
		connectionGroup.GET("/:id/schema", userHandler.List)

	}
}

func (s *Server) Start(addr string) error {
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
	}

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
