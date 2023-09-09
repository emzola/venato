package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/emzola/venato/pkg/discovery"
	"github.com/emzola/venato/pkg/discovery/consul"
	"github.com/emzola/venato/project/internal/controller/project"
	httpHandler "github.com/emzola/venato/project/internal/handler/http"
	"github.com/emzola/venato/project/internal/repository/postgresql"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var serviceName = "Project"

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	f, err := os.Open("base.yaml")
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
	h := httpHandler.New(ctrl)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), h.Routes()); err != nil {
		panic(err)
	}
}
