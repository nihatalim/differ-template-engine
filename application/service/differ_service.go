package service

import (
	"context"
	"differ-template-engine/application/domain"
	"differ-template-engine/application/repository"
	"differ-template-engine/application/request"
	"differ-template-engine/log"
	"differ-template-engine/pkg/client/nodiffer"
	"fmt"
	"strings"
)

type differService struct {
	templateRepository        repository.TemplateRepository
	executionResultRepository repository.ExecutionResultRepository
	nodifferApi               nodiffer.API
	logger                    log.Logger
}

var parameterTemplate = "${%s}"

func NewDifferService(templateRepository repository.TemplateRepository, executionResultRepository repository.ExecutionResultRepository, nodifferApi nodiffer.API, logger log.Logger) *differService {
	return &differService{templateRepository: templateRepository, executionResultRepository: executionResultRepository, nodifferApi: nodifferApi, logger: logger}
}

func (d *differService) HasDifference(userId string, req request.DifferExecutionRequest) (bool, error) {
	leftTemplate, err := d.templateRepository.GetTemplateByUserIdAndTemplateId(userId, req.Templates.Left.Id)
	if err != nil {
		return false, err
	}

	rightTemplate, err := d.templateRepository.GetTemplateByUserIdAndTemplateId(userId, req.Templates.Right.Id)
	if err != nil {
		return false, err
	}

	diffRequest := nodiffer.HasDiffRequest{
		OperationId: req.OperationId,
		Templates: nodiffer.Templates{
			Left: nodiffer.Template{
				Id:      leftTemplate.Id,
				Url:     NormalizeUrl(string(leftTemplate.Content.Url), req.Execution.Params),
				Method:  string(leftTemplate.Content.Method),
				Headers: NormalizeHeader(leftTemplate.Content.Headers, req.Execution.Params),
			},
			Right: nodiffer.Template{
				Id:      leftTemplate.Id,
				Url:     NormalizeUrl(string(rightTemplate.Content.Url), req.Execution.Params),
				Method:  string(leftTemplate.Content.Method),
				Headers: NormalizeHeader(rightTemplate.Content.Headers, req.Execution.Params),
			},
		},
	}

	diffResponse, err := d.nodifferApi.HasDiff(context.Background(), diffRequest)
	if err != nil {
		return false, err
	}

	if diffResponse.HasDiff {
		if err := d.executionResultRepository.SaveExecutionResult(diffRequest.OperationId, userId, diffRequest, diffResponse.HasDiff); err != nil {
			d.logger.Errorf("error while inserting execution result for operationId:%s", diffRequest.OperationId)
			return false, err
		}
	}

	return diffResponse.HasDiff, nil
}

func NormalizeUrl(url string, executionParameters map[string]string) string {
	for paramKey := range executionParameters {
		paramVal := executionParameters[paramKey]
		url = strings.ReplaceAll(url, fmt.Sprintf(parameterTemplate, paramKey), paramVal)
	}

	return url
}

func NormalizeHeader(templateHeaders domain.Headers, executionParameters map[string]string) map[string]string {
	headers := make(map[string]string)

	for templateHeaderKey := range templateHeaders {
		header := templateHeaders[templateHeaderKey]
		if header.Type == domain.StaticHeaderType {
			headers[templateHeaderKey] = *header.Value
			continue
		}

		if header.Type == domain.DynamicHeaderType {
			executionVal, ok := executionParameters[templateHeaderKey]
			if ok {
				headers[templateHeaderKey] = executionVal
			}
		}
	}

	return headers
}
