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
	CreatedBy      int       `gorm:"column:created_by;not null;default:0" json:"-"`
	UpdatedAt      time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"updatedAt"`
	UpdatedBy      int       `gorm:"column:updated_by;not null;default:0" json:"-"`
}
