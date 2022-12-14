package middleware

import (
	"golang-api/helper"
	"golang-api/service"

	"errors"
	"fmt"
	"net/http"
)

var (
	ErrTokenMalformed   = service.ErrTokenMalformed
	ErrTokenExpired     = service.ErrTokenExpired
	ErrTokenNotValidYet = service.ErrTokenNotValidYet
)

// AuthorizeJWT validates a JWT token, It also returns error is the token/encryption/issuer is invalid or sets a current user header
func AuthorizeJWT(service service.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response := helper.BuildErrorResponse(helper.COULD_NOT_PROCESS_REQUEST, errors.New("token not present in request"))

				helper.WriteJSON(w, http.StatusBadRequest, response)

				return
			}

			token, err := service.ValidateToken(authHeader)
			if !token.Valid {
				if errors.Is(err, ErrTokenMalformed) {
					response := helper.BuildErrorResponse(helper.ERROR_PROCESSING_REQUEST, errors.New("invalid token"))

					helper.WriteJSON(w, http.StatusUnauthorized, response)

					return
				} else if errors.Is(err, ErrTokenExpired) || errors.Is(err, ErrTokenNotValidYet) {
					response := helper.BuildErrorResponse(helper.ERROR_PROCESSING_REQUEST, errors.New("expired or inactive token"))

					helper.WriteJSON(w, http.StatusUnauthorized, response)

					return
				} else {
					response := helper.BuildErrorResponse(helper.COULD_NOT_PROCESS_REQUEST, errors.New("couldn't handle this token"))

					helper.WriteJSON(w, http.StatusUnauthorized, response)

					return
				}
			}

			claims := service.MapClaims(token)
			currentUserId := fmt.Sprintf("%v", claims["user_id"])
			tokenIssuer := fmt.Sprintf("%v", claims["iss"])

			if tokenIssuer != service.GetJWTIssuer() {
				response := helper.BuildErrorResponse(helper.ERROR_PROCESSING_REQUEST, errors.New("invalid token"))

				helper.WriteJSON(w, http.StatusBadRequest, response)

				return
			}

			r.Header.Add(helper.CURRENT_USER, currentUserId)

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

// ValidateUserInJWTProvided validates that AuthorizeJWT middleware has parsed and stored the authenticated user Id from the token passed into the request headers of the running context
func ValidateUserInJWTProvided(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authenticatedUserId := r.Header.Get(helper.CURRENT_USER)

		if authenticatedUserId == "" {
			response := helper.BuildErrorResponse(helper.ERROR_PROCESSING_REQUEST, errors.New("empty authentication context"))

			helper.WriteJSON(w, http.StatusBadRequest, response)

			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
