package authorization

import (
	"context"
	"log"

	"RedWood011/server/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MiddlewareJWT(cfg *config.Config) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if isNeedNoValidToken(info.FullMethod) {
			return handler(ctx, req)
		}
		// для всех остальных методов пропускаем запрос без проверки JWT токена
		userID, err := TokenValid(ctx, cfg.TokenConfig.SecretKey)
		log.Println("userID:", userID)
		if err != nil {

			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}
		ctx = context.WithValue(ctx, "userID", userID)
		return handler(ctx, req)
	}
}

func isNeedNoValidToken(method string) bool {
	return method == "/users.Users/CreateUser" || method == "/users.Users/AuthUser"
}
