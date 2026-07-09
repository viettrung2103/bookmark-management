package fixtures

import (
	"github.com/google/uuid"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	//"gorm.io/gorm"
)

func GetTestBase(baseUUID string) model.Base {
	return model.Base{
		ID: uuid.MustParse(baseUUID),
	}
}

func GetUUID(ObjectUUID string) uuid.UUID {
	return uuid.MustParse(ObjectUUID)
}
