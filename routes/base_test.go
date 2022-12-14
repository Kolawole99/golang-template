package routes

import (
	"golang-api/service"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func Test_HandleRouting(t *testing.T) {
	var db *gorm.DB
	var jwt service.JWTService

	testRoutes := HandleRouting(db, jwt)
	chiRoutes := testRoutes.(chi.Router)

	routes := []string{
		"/health",
		"/auth/register",
		"/auth/login",
		"/auth/forgot-password",
		"/auth/reset-password",
		"/users/update",
		"/users/profile",
	}

	for _, route := range routes {
		routeExists(t, chiRoutes, route)
	}
}

func routeExists(t *testing.T, routes chi.Router, route string) {
	found := false

	_ = chi.Walk(
		routes,
		func(
			method string,
			foundRoute string,
			handler http.Handler,
			middlewares ...func(http.Handler) http.Handler,
		) error {
			if "/api/v1"+route == foundRoute {
				found = true
			}
			return nil
		},
	)

	if !found {
		t.Errorf("did not find %s in registered routes", route)
	}
}
