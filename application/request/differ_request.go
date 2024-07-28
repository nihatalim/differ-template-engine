package request

type TemplateIdentification struct {
	Id int64 `json:"id"`
}

type TemplateIdentifications struct {
	Left  TemplateIdentification `json:"left"`
	Right TemplateIdentification `json:"right"`
}

type ExecutionParameters struct {
	Params map[string]string `json:"params"`
}

type DifferExecutionRequest struct {
	OperationId string                  `json:"operationId"`
	Templates   TemplateIdentifications `json:"templates"`
	Execution   ExecutionParameters     `json:"execution"`
}
