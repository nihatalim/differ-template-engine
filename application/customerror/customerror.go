package customerror

import "errors"

var (
	ErrTemplateNameIsNotValid = errors.New("template name is not valid!")
	ErrUrlIsNotValid          = errors.New("url is not valid!")
	ErrMethodIsNotValid       = errors.New("method is not valid!")
	ErrHeaderIsNotValid       = errors.New("header is not valid!")
	ErrTemplateIsNotFound     = errors.New("template is not found")
	ErrUserIdIsNotValid       = errors.New("userId is not valid")
	ErrTemplateIdIsNotValid   = errors.New("templateId is not valid")
	ErrUserIdIsRequired       = errors.New("x-userid is required, please pass valid user id in headers with x-userid key")
)
