package domain

import (
	"regexp"
)

type Url string

func (u *Url) GetParameters() []string {
	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(string(*u), -1)

	var params []string
	for _, match := range matches {
		if len(match) > 1 {
			params = append(params, match[1])
		}
	}
	return params
}

func (u *Url) Validate() bool {
	if u == nil {
		return false
	}

	if *u == "" {
		return false
	}

	return true
}
