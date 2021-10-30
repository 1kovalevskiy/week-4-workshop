package main

import (
	"context"
	"flag"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"

	grpc_category_service "github.com/ozonmp/week-3-workshop/category-service/pkg/category-service"
	"github.com/ozonmp/week-3-workshop/product-service/internal/config"
	"github.com/ozonmp/week-3-workshop/product-service/internal/server"
	product_service "github.com/ozonmp/week-3-workshop/product-service/internal/service/product"
)

func main() {
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}
	cfg := config.GetConfigInstance()

	flag.Parse()

	log.Info().
		Str("version", cfg.Project.Version).
		Str("commitHash", cfg.Project.CommitHash).
		Bool("debug", cfg.Project.Debug).
		Str("environment", cfg.Project.Environment).
		Msgf("Starting service: %s", cfg.Project.Name)

	// default: zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.Project.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	categoryServiceConn, err := grpc.DialContext(
		context.Background(),
		cfg.CategoryServiceAddr,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create client")
	}

	categoryServiceClient := grpc_category_service.NewCategoryServiceClient(categoryServiceConn)

	productService := product_service.NewService(categoryServiceClient)

	if err := server.NewGrpcServer(productService).Start(&cfg); err != nil {
		log.Error().Err(err).Msg("Failed creating gRPC server")

		return
	}
}
