package controller

import (
	"differ-template-engine/application/customerror"
	"differ-template-engine/application/domain"
	"differ-template-engine/application/request"
	"differ-template-engine/log"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type UserController struct {
	logger  log.Logger
	service UserService
}

type UserService interface {
	GetUserTemplates(userId string) []domain.Template
	CreateTemplate(userId string, template domain.Template) (*domain.Template, error)
	DeleteTemplate(userId string, templateId int64) error
}

func NewUserController(service UserService, logger log.Logger) *UserController {
	return &UserController{
		logger:  logger,
		service: service,
	}
}

const EndpointGetTemplatesByUserId = "/users/:userId/templates"
const EndpointDeleteTemplateByUserIdAndTemplateId = "/users/:userId/templates/:templateId"
const EndpointCreateTemplateByUserId = "/users/:userId/templates"

func (u *UserController) RegisterRoutes(f *fiber.App) {
	f.Get(EndpointGetTemplatesByUserId, u.ProcessGetTemplatesByUserId)
	f.Delete(EndpointDeleteTemplateByUserIdAndTemplateId, u.ProcessDeleteTemplateByUserIdAndTemplateId)
	f.Post(EndpointCreateTemplateByUserId, u.ProcessCreateTemplateByUserId)
}

// @Tags controller
// @Accept json
// @Produce json
// @Success 200 {object} []domain.Template
// @Param 	userId 	path	string	true "userId"
// @Router /users/{userId}/templates [get]
func (u *UserController) ProcessGetTemplatesByUserId(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")
	if userId == "" {
		return customerror.ErrUserIdIsNotValid
	}

	templates := u.service.GetUserTemplates(userId)
	resp, _ := sonic.Marshal(templates)

	return ctx.Status(fiber.StatusOK).Send(resp)
}

// @Tags controller
// @Accept json
// @Produce json
// @Success 204
// @Param 	userId 		path 	string 	true 	"userId"
// @Param 	templateId 	path 	string 	true 	"templateId"
// @Router /users/{userId}/templates/{templateId} [delete]
func (u *UserController) ProcessDeleteTemplateByUserIdAndTemplateId(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")
	if userId == "" {
		return customerror.ErrUserIdIsNotValid
	}

	templateId, err := strconv.ParseInt(ctx.Params("templateId"), 10, 64)
	if err != nil {
		return customerror.ErrTemplateIdIsNotValid
	}

	err = u.service.DeleteTemplate(userId, templateId)
	if err != nil {
		if err == customerror.ErrTemplateIsNotFound {
			return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
		}

		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

// @Tags controller
// @Accept json
// @Produce json
// @Success 201 {object} domain.Template
// @Param 	requestBody 	body	request.CreateTemplateRequest 	true 	"CreateTemplateRequest"
// @Param 	userId 	path	string	true "userId"
// @Router /users/{userId}/templates [post]
func (u *UserController) ProcessCreateTemplateByUserId(ctx *fiber.Ctx) error {
	var req request.CreateTemplateRequest
	body := ctx.Body()

	if err := sonic.Unmarshal(body, &req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := req.Validate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	userId := ctx.Params("userId")
	if userId == "" {
		return customerror.ErrUserIdIsNotValid
	}

	template, err := u.service.CreateTemplate(userId, req.ToTemplate())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	resp, _ := sonic.Marshal(template)

	return ctx.Status(fiber.StatusCreated).Send(resp)
}
