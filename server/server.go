package server

import (
	"differ-template-engine/log"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/swagger"
)

type Server struct {
	app    *fiber.App
	logger log.Logger
}

type RegisterRoutesFunc func(app *fiber.App)

func New(logger log.Logger, registerRoutesFunc []RegisterRoutesFunc) Server {
	fiberApp := fiber.New(
		fiber.Config{
			ReadBufferSize:           8192,
			DisableDefaultDate:       true,
			DisableHeaderNormalizing: true,
			JSONEncoder:              sonic.Marshal,
			JSONDecoder:              sonic.Unmarshal,
		},
	)

	fiberApp.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	fiberApp.Use(pprof.New())

	fiberApp.Get("/swagger/*", swagger.HandlerDefault)

	fiberApp.Use(func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "application/json")
		return ctx.Next()
	})

	for _, registerRouteFunc := range registerRoutesFunc {
		registerRouteFunc(fiberApp)
	}

	server := Server{logger: logger, app: fiberApp}

	return server
}

func (s Server) Run(port string) error {
	s.logger.Infof("Differ template engine is listening %s", port)
	return s.app.Listen(port)
}

func (s Server) Close() error {
	s.logger.Info("Differ template engine is closing...")
	return s.app.Shutdown()
}
