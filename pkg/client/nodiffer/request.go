package nodiffer

type Template struct {
	Id      int64             `json:"id"`
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
}

type Templates struct {
	Left  Template `json:"left"`
	Right Template `json:"right"`
}

type HasDiffRequest struct {
	OperationId string    `json:"operationId"`
	Templates   Templates `json:"templates"`
}
