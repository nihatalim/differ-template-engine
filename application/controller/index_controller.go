package controller

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

const (
	EndpointIndex = "/"
)

type IndexController struct{}

func NewIndexController() *IndexController {
	return &IndexController{}
}

func (c *IndexController) RegisterRoutes(f *fiber.App) {
	f.Get(EndpointIndex, c.IndexHandler)
}

func (c *IndexController) IndexHandler(ctx *fiber.Ctx) error {
	return ctx.Redirect("/swagger/index.html", http.StatusPermanentRedirect)
}
