package routes

import (
	"errors"
	_ "golang-api/docs"
	"golang-api/helper"
	"golang-api/service"

	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

type Routing struct {
	db         *gorm.DB
	jwtService service.JWTService
}

// HandleRouting enables application wide routing
func HandleRouting(db *gorm.DB, jwt service.JWTService) http.Handler {
	routeService := Routing{
		jwtService: jwt,
		db:         db,
	}

	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))
	mux.NotFound(notFound)

	mux.Route("/api/v1", func(routes chi.Router) {
		routes.Get("/health", healthCheck)

		authRoutes(routes, &routeService)
		userRoutes(routes, &routeService)
	})

	mux.Mount("/api/v1/docs", httpSwagger.WrapHandler)

	return mux
}

// healthCheck godoc
// @Summary     Shows the health status of server.
// @Description This ping route enables us verify the health status and availability of the server.
// @Tags        Health Check
// @Accept      */*
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Router      /health [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	helper.WriteJSON(w, http.StatusOK, helper.Response{
		Message: "Application is up, running and healthy!",
	})
}

// notFound is the custom response for all 404 request in the application
func notFound(w http.ResponseWriter, r *http.Request) {
	response := helper.BuildErrorResponse("Not Found", errors.New("not found"))

	helper.WriteJSON(w, http.StatusNotFound, response)
}
