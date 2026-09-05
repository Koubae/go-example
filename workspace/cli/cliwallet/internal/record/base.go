package record

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UUIDModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (m *UUIDModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
