package directive

import (
	"github.com/satori/go.uuid"
)

type BaseDirective struct {
	Type string `json:"type"`
}

func (this *BaseDirective) GenToken() string {
	return uuid.NewV4().String()
}
