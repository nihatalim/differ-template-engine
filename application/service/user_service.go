package service

import (
	"differ-template-engine/application/domain"
	"differ-template-engine/application/repository"
	"differ-template-engine/log"
)

type userService struct {
	repo   repository.TemplateRepository
	logger log.Logger
}

func NewUserService(repo repository.TemplateRepository, logger log.Logger) *userService {
	return &userService{repo: repo, logger: logger}
}

func (u *userService) GetUserTemplates(userId string) []domain.Template {
	templates, err := u.repo.GetUserTemplate(userId)
	if err != nil {
		u.logger.Errorf("failed to fetch user templates %s", userId)
	}

	return templates
}

func (u *userService) CreateTemplate(userId string, t domain.Template) (*domain.Template, error) {
	template, err := u.repo.SaveTemplate(userId, t)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (u *userService) DeleteTemplate(userId string, templateId int64) error {
	if err := u.repo.DeleteTemplate(userId, templateId); err != nil {
		u.logger.Errorf("failed to delete user template, userid: %s templateid: ", userId, templateId)
		return err
	}

	return nil
}
