package controller

import (
	"differ-template-engine/application/customerror"
	"differ-template-engine/application/request"
	"differ-template-engine/log"
	"differ-template-engine/response"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
)

const EndpointDifferExecute = "/differ/execute"

type DifferController struct {
	logger  log.Logger
	service DifferService
}

type DifferService interface {
	HasDifference(userId string, request request.DifferExecutionRequest) (bool, error)
}

func NewDifferController(service DifferService, logger log.Logger) *DifferController {
	return &DifferController{
		logger:  logger,
		service: service,
	}
}

func (d *DifferController) RegisterRoutes(f *fiber.App) {
	f.Post(EndpointDifferExecute, d.ProcessDifferExecute)
}

func getUserId(ctx fiber.Ctx) string {
	headers := ctx.GetReqHeaders()
	userIds := headers["x-userid"]
	if len(userIds) == 0 {
		return ""
	}

	return userIds[0]
}

func (d *DifferController) ProcessDifferExecute(ctx fiber.Ctx) error {
	var req request.DifferExecutionRequest
	userId := getUserId(ctx)
	if userId == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString(customerror.ErrUserIdIsRequired.Error())
	}

	body := ctx.Body()

	if err := sonic.Unmarshal(body, &req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	hasDifference, err := d.service.HasDifference(userId, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	executionResponse := response.DifferExecutionResponse{HasDifference: hasDifference}
	responseBody, _ := sonic.Marshal(executionResponse)

	return ctx.Status(fiber.StatusOK).Send(responseBody)
}
