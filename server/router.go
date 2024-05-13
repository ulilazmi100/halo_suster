package server

import (
	"halo_suster/controller"
	"halo_suster/middleware"
	"halo_suster/repo"
	"halo_suster/svc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func (s *Server) RegisterRoute() {
	mainRoute := s.app.Group("/v1")

	registerUserRoute(mainRoute, s.dbPool)
}

func registerUserRoute(r *echo.Group, db *pgxpool.Pool) {
	ctr := controller.NewUserController(svc.NewUserSvc(repo.NewUserRepo(db)))
	itGroup := r.Group("/user/it")

	newRoute(itGroup, "POST", "/register", ctr.Register)
	newRoute(itGroup, "POST", "/login", ctr.Login)
}

func newRoute(router *echo.Group, method, path string, handler echo.HandlerFunc) {
	router.Add(method, path, handler)
}

func newRouteWithAuth(router *echo.Group, method, path string, handler echo.HandlerFunc) {
	router.Add(method, path, handler, middleware.Auth)
}
