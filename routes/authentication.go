package routes

import (
	"golang-api/controller"
	"golang-api/model"
	"golang-api/service"

	"github.com/go-chi/chi/v5"
)

func authRoutes(appRouter chi.Router, routesService *Routing) {
	var (
		userRepository model.UserRepository = model.NewUserRepository(routesService.db)

		authService service.AuthService = service.NewAuthService(userRepository, routesService.jwtService)

		authController controller.AuthController = controller.NewAuthController(authService)
	)

	appRouter.Route("/auth", func(routes chi.Router) {
		routes.Post("/login", authController.Login)
		routes.Post("/register", authController.Register)
		routes.Post("/forgot-password", authController.ForgotPassword)
		routes.Post("/reset-password", authController.ResetPassword)
	})
}
