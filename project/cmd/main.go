package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/emzola/venato/gen"
	"github.com/emzola/venato/pkg/discovery"
	"github.com/emzola/venato/pkg/discovery/consul"
	"github.com/emzola/venato/project/internal/controller/project"
	grpcHandler "github.com/emzola/venato/project/internal/handler/grpc"
	"github.com/emzola/venato/project/internal/repository/postgresql"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

var serviceName = "Project"

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	f, err := os.Open(filepath.FromSlash("../configs/base.yaml"))
	if err != nil {
		logger.Fatal("Failed to open configuration", zap.Error(err))
	}
	var cfg config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		logger.Fatal("Failed to parse configuration", zap.Error(err))
	}
	port := cfg.API.Port
	logger.Info("Starting the project service", zap.Int("port", port))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				logger.Error("Failed to report healthy state", zap.Error(err))
			}
			time.Sleep(time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)
	repo, err := postgresql.New()
	if err != nil {
		logger.Fatal("Failed to establish database connection pool", zap.Error(err))
	}
	ctrl := project.New(repo)
	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterProjectServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
