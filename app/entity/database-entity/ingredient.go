package databaseentity

import "time"

// table ingredients model
type Ingredient struct {
	ID           int       `gorm:"primaryKey" json:"-"`
	Serial       string    `gorm:"unique;size:11;not null" json:"serial"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	UnitID       int       `gorm:"not null" json:"-"`
	Unit         *Unit     `gorm:"foreignKey:unit_id" json:"unit"`
	PricePerUnit float32   `gorm:"not null;type:decimal(10,2);default:0" json:"pricePerUnit"`
	Stock        float32   `gorm:"not null;type:decimal(10,2);default:0" json:"stock"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	CreatedByDB  int       `gorm:"column:created_by;not null;default:0" json:"-"`
	CreatedBy    string    `gorm:"-" json:"createdBy"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"updatedAt"`
	UpdatedByDB  int       `gorm:"column:updated_by;not null;default:0" json:"-"`
	UpdatedBy    string    `gorm:"-" json:"updatedBy"`
}
