package server

import (
	configs "halo_suster/cfg"
	"halo_suster/controller"
	"halo_suster/middleware"
	"halo_suster/repo"
	"halo_suster/svc"
	"context"

	"log"

	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func (s *Server) RegisterRoute(config configs.Config) {
	mainRoute := s.app.Group("/v1")

	registerUserRoute(mainRoute, s.db)
	registerMedicalRoute(mainRoute, s.db)
	registerImageRoute(mainRoute, config)
}

func registerUserRoute(r *echo.Group, db *sqlx.DB) {
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
	newRouteWithItAuth(userGroup, "GET", "", ctr.GetUser)
}

func registerMedicalRoute(r *echo.Group, db *sqlx.DB) {
	ctr := controller.NewMedicalController(svc.NewMedicalSvc(repo.NewMedicalRepo(db)))

	patientGroup := r.Group("/medical/patient")
	newRouteWithAuth(patientGroup, "POST", "", ctr.RegisterPatient)
	newRouteWithAuth(patientGroup, "GET", "", ctr.GetPatient)

	recordGroup := r.Group("/medical/record")
	newRouteWithAuth(recordGroup, "POST", "", ctr.RegisterRecord)
	newRouteWithAuth(recordGroup, "GET", "", ctr.GetRecord)
}

func registerImageRoute(r *echo.Group, config configs.Config) {
	bucket := config.AWS_S3_BUCKET_NAME

	// Load AWS configuration
	cfg, err := awsCfg.LoadDefaultConfig(
		context.Background(),
		awsCfg.WithRegion("ap-southeast-1"),
		awsCfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.AWS_ACCESS_KEY_ID,
				config.AWS_SECRET_ACCESS_KEY,
				"",
			),
		),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	// Initialize the image service and controller
	imageService := svc.NewImageSvc(cfg, bucket)
	imageController := controller.NewImageController(imageService)

	// Register the route with authentication middleware
	newRouteWithAuth(r, "POST", "/image", imageController.UploadImage)
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
