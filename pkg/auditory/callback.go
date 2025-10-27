package auditory

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
)

func GeneratePersistentEvents(db *gorm.DB) (*EntityEvent, error) {
	if db.Error != nil {
		return nil, db.Error
	}

	obj := db.Statement.Dest

	auditable, ok := obj.(Auditable)
	if !ok {
		return nil, fmt.Errorf("not a auditable object")
	}

	contextID := db.Statement.Context.Value("context_id")
	userId := db.Statement.Context.Value("user_id")
	userpool := db.Statement.Context.Value("userpool")
	if contextID == nil {
		logrus.Warnf("Nao foi possivel obter o 'context_id' do contexto")
	}

	if userId == nil {
		logrus.Warnf("Nao foi possivel obter o 'user_id' do contexto")
	}

	if userpool == nil {
		logrus.Warnf("Nao foi possivel obter o 'userpool' do contexto")
	}

	fields := db.Statement.ReflectValue.NumField()

	entityFields := make([]EntityField, fields)

	for i := 0; i < fields; i++ {
		field := db.Statement.ReflectValue.Type().Field(i)
		gormTag := field.Tag.Get("gorm")
		entityFields[i] = EntityField{
			Name:     db.Statement.ReflectValue.Type().Field(i).Name,
			GormTags: strings.Split(gormTag, ";"),
		}
	}

	content, err := json.Marshal(db.Statement.Dest)

	if err != nil {
		logrus.Errorf("error marshalling entity: %v", err)
	}

	return &EntityEvent{
		Key:       auditable.GetKey(),
		Table:     db.Statement.Table,
		Content:   base64.StdEncoding.EncodeToString(content),
		ContextID: uuid.Must(uuid.FromString(fmt.Sprintf("%v", contextID))),
		UserID:    uuid.Must(uuid.FromString(fmt.Sprintf("%v", userId))),
		Fields:    entityFields,
		UserPool:  fmt.Sprintf("%v", userpool),
	}, nil

}
