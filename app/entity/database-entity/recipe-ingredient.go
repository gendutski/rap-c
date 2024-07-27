package databaseentity

import "time"

// table recipe_ingredients model
type RecipeIngredient struct {
	ID           int         `gorm:"primaryKey" json:"-"`
	Serial       string      `gorm:"unique;size:11;not null" json:"serial"`
	RecipeID     int         `gorm:"not null" json:"-"`
	Recipe       *Recipe     `gorm:"foreignKey:recipe_id" json:"-"`
	IngredientID int         `gorm:"not null" json:"-"`
	Ingredient   *Ingredient `gorm:"foreignKey:ingredient_id" json:"ingredient"`
	UnitID       int         `gorm:"not null" json:"-"`
	Unit         *Unit       `gorm:"foreignKey:unit_id" json:"unit"`
	Quantity     float32     `gorm:"not null;type:decimal(10,2);default:0" json:"quantity"`
	CreatedAt    time.Time   `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	CreatedByDB  int         `gorm:"column:created_by;not null;default:0" json:"-"`
	CreatedBy    string      `gorm:"-" json:"createdBy"`
	UpdatedAt    time.Time   `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"updatedAt"`
	UpdatedByDB  int         `gorm:"column:updated_by;not null;default:0" json:"-"`
	UpdatedBy    string      `gorm:"-" json:"updatedBy"`
}
