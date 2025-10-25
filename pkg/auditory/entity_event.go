package auditory

import "github.com/gofrs/uuid"

type EntityEvent struct {
	Key       string        `json:"key"`
	ContextID uuid.UUID     `json:"context_id"`
	UserID    uuid.UUID     `json:"user_id"`
	Table     string        `json:"table"`
	Content   string        `json:"content"`
	Fields    []EntityField `json:"fields"`
	UserPool  string        `json:"user_pool"`
}
type EntityField struct {
	Name     string   `json:"name"`
	GormTags []string `json:"gorm_tags"`
}

type Auditable interface {
	GetKey() string
}
