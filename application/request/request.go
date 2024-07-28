package request

import (
	"differ-template-engine/application/customerror"
	"differ-template-engine/application/domain"
)

type CreateTemplateRequest struct {
	Name    string         `json:"name"`
	Url     domain.Url     `json:"url"`
	Method  domain.Method  `json:"method"`
	Headers domain.Headers `json:"headers"`
}

type DeleteTemplateRequest struct{}

func (ctr *CreateTemplateRequest) Validate() error {
	if ctr.Name == "" {
		return customerror.ErrTemplateNameIsNotValid
	}

	if !ctr.Url.Validate() {
		return customerror.ErrUrlIsNotValid
	}

	if !ctr.Method.Validate() {
		return customerror.ErrMethodIsNotValid
	}

	if !ctr.Headers.Validate() {
		return customerror.ErrHeaderIsNotValid
	}

	return nil
}

func (ctr *CreateTemplateRequest) ToTemplate() domain.Template {
	return domain.Template{
		Name: ctr.Name,
		Content: domain.TemplateContent{
			Url:     ctr.Url,
			Method:  ctr.Method,
			Headers: ctr.Headers,
		},
	}
}
