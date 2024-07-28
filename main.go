package main

import (
	"context"
	"differ-template-engine/application/controller"
	"differ-template-engine/application/repository"
	"differ-template-engine/application/service"
	_ "differ-template-engine/docs"
	"differ-template-engine/log"
	"differ-template-engine/pkg/client/nodiffer"
	"differ-template-engine/pkg/config"
	"differ-template-engine/server"
	"go.uber.org/zap/zapcore"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := log.New(zapcore.InfoLevel)

	cfg, err := config.New()
	if err != nil {
		panic("application configuration could not read, error " + err.Error())
	}

	conn, err := repository.Initialize(cfg, logger)
	if err != nil {
		panic("db connection could not initialize, error " + err.Error())
	}

	defer conn.Close(ctx)

	templateRepository := repository.NewTemplateRepository(cfg, conn, logger)
	executionResultRepository := repository.NewExecutionResultRepository(cfg, conn, logger)
	userService := service.NewUserService(templateRepository, logger)
	userController := controller.NewUserController(userService, logger)

	nodifferApi := nodiffer.NewNodifferAPI(cfg.Clients.Nodiffer)
	differService := service.NewDifferService(templateRepository, executionResultRepository, nodifferApi, logger)
	differController := controller.NewDifferController(differService, logger)

	indexController := controller.NewIndexController()

	svr := server.New(logger, []server.RegisterRoutesFunc{
		indexController.RegisterRoutes,
		userController.RegisterRoutes,
		differController.RegisterRoutes,
	})

	go func() {
		if err := svr.Run(cfg.Port); err != nil {
			logger.Errorf("Server run error %s", err.Error())
		}
	}()

	<-ctx.Done()
	_ = svr.Close()
}
