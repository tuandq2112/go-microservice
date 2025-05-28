package interceptors

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tuandq2112/go-microservices/shared/errors"
	"github.com/tuandq2112/go-microservices/shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type ContextKey string

const (
	UserClaimsKey ContextKey = "user_claims"
)

func UnaryAuthInterceptor(logger logger.Logger, jwtSecret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == "/user.UserService/Login" || info.FullMethod == "/user.UserService/Register" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.UnauthorizedError()
		}

		authHeader := md.Get("grpcgateway-authorization")
		if len(authHeader) == 0 {
			return nil, errors.UnauthorizedError()
		}

		tokenString := authHeader[0]
		claims, err := validateToken(tokenString, jwtSecret)
		if err != nil {
			logger.Error("Invalid token", zap.Error(err))
			return nil, errors.UnauthorizedError()
		}

		newCtx := context.WithValue(ctx, UserClaimsKey, claims)
		return handler(newCtx, req)
	}
}

func StreamAuthInterceptor(logger logger.Logger, jwtSecret string) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if info.FullMethod == "/user.UserService/Login" || info.FullMethod == "/user.UserService/Register" {
			return handler(srv, ss)
		}

		ctx := ss.Context()

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return errors.UnauthorizedError()
		}
		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return errors.UnauthorizedError()
		}

		tokenString := authHeader[0]
		claims, err := validateToken(tokenString, jwtSecret)
		if err != nil {
			logger.Error("Invalid token", zap.Error(err))
			return errors.UnauthorizedError()
		}

		newCtx := context.WithValue(ctx, UserClaimsKey, claims)
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          newCtx,
		}

		return handler(srv, wrappedStream)
	}
}

func validateToken(tokenString string, jwtSecret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
