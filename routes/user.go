package routes

import (
	"golang-api/controller"
	"golang-api/middleware"
	"golang-api/model"
	"golang-api/service"

	"github.com/go-chi/chi/v5"
)

func userRoutes(appRouter chi.Router, routesService *Routing) {
	var (
		userRepository model.UserRepository = model.NewUserRepository(routesService.db)

		userService service.UserService = service.NewUserService(userRepository)

		userController controller.UserController = controller.NewUserController(userService)
	)

	appRouter.Route("/users", func(routes chi.Router) {
		routes.Use(middleware.AuthorizeJWT(routesService.jwtService))
		routes.Use(middleware.ValidateUserInJWTProvided)

		routes.Put("/update", userController.UpdateUser)
		routes.Get("/profile", userController.Profile)
	})
}
