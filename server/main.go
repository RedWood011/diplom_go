package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"RedWood011/server/internal/authorization"
	"RedWood011/server/internal/config"
	"RedWood011/server/internal/database/postgres"
	"RedWood011/server/internal/logger"
	"RedWood011/server/internal/services/secret"
	"RedWood011/server/internal/services/user"
	secretgrpc "RedWood011/server/internal/transport/grpc/secret"
	usergrpc "RedWood011/server/internal/transport/grpc/user"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := logger.InitLogger()
	cfg := config.NewConfig()
	repo, err := postgres.NewDatabase(ctx, cfg.Database.URI, cfg.Database.MaxAttempts)
	if err != nil {
		logger.Info("Error database:", err.Error())
		log.Fatal(err)

	}
	err = repo.Ping(ctx)
	if err != nil {
		logger.Info("Error ping database:", err.Error())
		log.Fatal(err)
	}

	userService := user.NewUserService(repo, logger)
	grpcUsers := usergrpc.NewGrpcUsers(userService, cfg, logger)
	secretService := secret.NewSecretService(repo, logger)
	grpcSecrets := secretgrpc.NewGrpcSecrets(secretService, cfg, logger)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(authorization.MiddlewareJWT(cfg))),
	)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
		if err != nil {
			logger.Info("gRPC server failed to listen:", err.Error())
			return err
		}
		logger.Info("gRPC server listening:", lis.Addr())
		usergrpc.RegisterUsersServer(grpcServer, grpcUsers)
		secretgrpc.RegisterSecretsServer(grpcServer, grpcSecrets)
		return grpcServer.Serve(lis)
	})

	err = g.Wait()
	if err != nil {
		slog.Info("server returning an error:", err.Error())
	}

}
