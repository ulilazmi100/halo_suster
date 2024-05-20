package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	configs "halo_suster/cfg"
	local_mid "halo_suster/middleware"
)

type Server struct {
	dbPool    *pgxpool.Pool
	app       *echo.Echo
	validator *validator.Validate
}

func NewServer(db *pgxpool.Pool) *Server {
	// Create an Echo instance
	app := echo.New()
	app.HTTPErrorHandler = local_mid.ErrorHandler

	// Initialize validator
	validate := validator.New()

	// Middleware
	app.Use(middleware.Recover())
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())
	app.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	return &Server{
		dbPool:    db,
		app:       app,
		validator: validate,
	}
}

func (s *Server) Run(config configs.Config) {
	s.app.Logger.Fatal(s.app.Start(":" + config.APPPort))
}
