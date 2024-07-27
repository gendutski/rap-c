package databaseentity

import "time"

// table products model
type Product struct {
	ID             int       `gorm:"primaryKey" json:"-"`
	Serial         string    `gorm:"unique;size:11;not null" json:"serial"`
	RecipeID       int       `gorm:"not null" json:"-"`
	Recipe         *Recipe   `gorm:"foreignKey:recipe_id" json:"recipe"`
	Date           time.Time `gorm:"type:date;not null" json:"date"`
	Quantity       int       `gorm:"not null;default:0" json:"quantity"`
	SoldQuantity   int       `gorm:"not null;default:0" json:"soldQuantity"`
	ProfitExpected float32   `gorm:"not null;type:decimal(10,2);default:0" json:"profitExpected"`
	ProfitGet      float32   `gorm:"not null;type:decimal(10,2);default:0" json:"profitGet"`
	Status         string    `gorm:"type:enum('in production', 'in sales', 'sent to journal');not null;index" json:"status"`
	CreatedAt      time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	CreatedByDB    int       `gorm:"column:created_by;not null;default:0" json:"-"`
	CreatedBy      string    `gorm:"-" json:"createdBy"`
	UpdatedAt      time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"updatedAt"`
	UpdatedByDB    int       `gorm:"column:updated_by;not null;default:0" json:"-"`
	UpdatedBy      string    `gorm:"-" json:"updatedBy"`
}
