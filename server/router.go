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
	registerMedicalRoute(mainRoute, s.dbPool)
}

func registerUserRoute(r *echo.Group, db *pgxpool.Pool) {
	ctr := controller.NewUserController(svc.NewUserSvc(repo.NewUserRepo(db)))

	itGroup := r.Group("/user/it")
	newRoute(itGroup, "POST", "/register", ctr.Register)
	newRoute(itGroup, "POST", "/login", ctr.Login)

	nurseGroup := r.Group("/user/nurse")
	newRouteWithItAuth(nurseGroup, "POST", "/register", ctr.NurseRegister)
	newRoute(nurseGroup, "POST", "/login", ctr.NurseLogin)
	newRouteWithItAuth(nurseGroup, "PUT", "/:userId", ctr.NurseUpdate)
	newRouteWithItAuth(nurseGroup, "DELETE", "/:userId", ctr.NurseDelete)
	newRouteWithItAuth(nurseGroup, "POST", "/:userId/access", ctr.NurseAccess)

	userGroup := r.Group("/user")
	newRouteWithItAuth(userGroup, "GET", "/", ctr.GetUser)
}

func registerPatientRoute(r *echo.Group, db *pgxpool.Pool) {
	ctr := controller.NewPatientController(svc.NewPatientSvc(repo.NewPatientRepo(db)))
	patientGroup := r.Group("/medical/patient")

	newRouteWithAuth(patientGroup, "POST", "/", ctr.RegisterPatient)
	newRouteWithAuth(patientGroup, "GET", "/", ctr.GetPatient)
}

func newRoute(router *echo.Group, method, path string, handler echo.HandlerFunc) {
	router.Add(method, path, handler)
}

func newRouteWithAuth(router *echo.Group, method, path string, handler echo.HandlerFunc) {
	router.Add(method, path, handler, middleware.Auth)
}

func newRouteWithItAuth(router *echo.Group, method, path string, handler echo.HandlerFunc) {
	router.Add(method, path, handler, middleware.ItAuth)
}
