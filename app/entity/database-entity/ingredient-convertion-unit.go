package databaseentity

import "time"

// table ingredient_conversion_units model
type IngredientConvertionUnit struct {
	ID            int         `gorm:"primaryKey" json:"-"`
	Serial        string      `gorm:"unique;size:11;not null" json:"serial"`
	IngredientID  int         `gorm:"not null" json:"-"`
	Ingredient    *Ingredient `gorm:"foreignKey:ingredient_id" json:"ingredient"`
	UnitID        int         `gorm:"not null" json:"-"`
	Unit          *Unit       `gorm:"foreignKey:unit_id" json:"unit"`
	Value         float32     `gorm:"not null;type:decimal(10,2);default:0" json:"value"`
	SkipCalculate bool        `gorm:"not null;default:0" json:"skipCalculate"`
	CreatedAt     time.Time   `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	CreatedBy     int         `gorm:"column:created_by;not null;default:0" json:"-"`
	UpdatedAt     time.Time   `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"updatedAt"`
	UpdatedBy     int         `gorm:"column:updated_by;not null;default:0" json:"-"`
}
