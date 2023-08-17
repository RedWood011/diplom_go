package authorization

import (
	"context"

	"RedWood011/server/internal/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGUID string
type UserKey string

func (uk UserKey) String() string {
	return string(uk)
}
func (ug UserGUID) String() string {
	return string(ug)
}

func MiddlewareJWT(cfg *config.Config) func(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if isNeedNoValidToken(info.FullMethod) {
			return handler(ctx, req)
		}

		userID, err := TokenValid(ctx, cfg.TokenConfig.SecretKey)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}
		ctx = context.WithValue(ctx, UserKey("userID"), UserGUID(userID))
		return handler(ctx, req)
	}
}

func isNeedNoValidToken(method string) bool {
	return method == "/users.Users/CreateUser" || method == "/users.Users/AuthUser"
}
