package domain

import "time"

type TemplateContent struct {
	Url     Url     `json:"url"`
	Method  Method  `json:"method"`
	Headers Headers `json:"headers"`
}

type Template struct {
	Id          int64           `json:"id"`
	UserId      string          `json:"userId"`
	Name        string          `json:"name"`
	Content     TemplateContent `json:"content"`
	CreatedDate time.Time       `json:"created_date"`
}
