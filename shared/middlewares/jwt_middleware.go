package middlewares

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/tuandq2112/go-microservices/shared/errors"
)

type Whitelist struct {
	Path   string
	Method string
}

type WhitelistPaths struct {
	Get    []string
	Post   []string
	Put    []string
	Delete []string
}

func routeToRegex(route string) string {
	escaped := regexp.QuoteMeta(route)

	escaped = strings.ReplaceAll(escaped, `\*\*`, `.*`)

	escaped = strings.ReplaceAll(escaped, `\*`, `[^/]+`)

	return "^" + escaped + "$"
}

func isRouteMatch(route, path string) bool {
	pattern := routeToRegex(route)
	matched, err := regexp.MatchString(pattern, path)
	if err != nil {
		return false
	}
	return matched
}

func isWhitelisted(path, method string, whitelistPaths WhitelistPaths) bool {
	if strings.HasPrefix(path, "/swagger/") {
		return true
	}

	var whitelist []string
	switch method {
	case "GET":
		whitelist = whitelistPaths.Get
	case "POST":
		whitelist = whitelistPaths.Post
	case "PUT":
		whitelist = whitelistPaths.Put
	case "DELETE":
		whitelist = whitelistPaths.Delete
	default:
		return false
	}

	for _, route := range whitelist {
		if isRouteMatch(route, path) {
			return true
		}
	}
	return false
}

func JWTMiddleware(whitelistPaths WhitelistPaths, userContextKey interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isWhitelisted(r.URL.Path, r.Method, whitelistPaths) {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				localizedMessage := T(r, errors.ErrUnauthorized, nil)
				errors.WriteError(w, errors.UnauthorizedError, localizedMessage, nil)
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				localizedMessage := T(r, errors.ErrBadRequest, map[string]interface{}{
					"details": "Invalid Authorization header format",
				})
				errors.WriteError(w, errors.BadRequestError, localizedMessage, nil)
				return
			}

			tokenString := tokenParts[1]
			claims, err := ParseJWT(tokenString)
			if err != nil {
				// Handle connection errors
				if errors.IsConnectionError(err) {
					localizedMessage := T(r, errors.ErrConnectionFailed, map[string]interface{}{
						"details": "Unable to connect to authentication service",
					})
					errors.WriteError(w, errors.ConnectionError, localizedMessage, err)
					return
				}

				// Handle internal errors (like JWT parsing errors)
				if strings.Contains(err.Error(), "internal") {
					localizedMessage := T(r, errors.ErrInternalServer, map[string]interface{}{
						"details": "Error processing authentication token",
					})
					errors.WriteError(w, errors.InternalServerError, localizedMessage, err)
					return
				}

				// Handle invalid/expired token
				localizedMessage := T(r, errors.ErrUnauthorized, map[string]interface{}{
					"details": "Invalid or expired token",
				})
				errors.WriteError(w, errors.UnauthorizedError, localizedMessage, err)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
