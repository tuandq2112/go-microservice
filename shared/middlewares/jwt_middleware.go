package middlewares

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/tuandq2112/go-microservices/shared/errors"
	"google.golang.org/grpc/codes"
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
				err := errors.UnauthorizedError()
				writeError(w, err)
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				err := errors.BadRequestError()
				writeError(w, err)
				return
			}

			tokenString := tokenParts[1]
			claims, err := ParseJWT(tokenString)
			if err != nil {
				err := errors.UnauthorizedError()
				writeError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeError(w http.ResponseWriter, err error) {
	if customErr, ok := err.(interface{ Code() codes.Code }); ok {
		w.WriteHeader(int(customErr.Code()))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(err.Error()))
}
