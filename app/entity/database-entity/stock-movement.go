package databaseentity

import "time"

// table stock_movements model
type StockMovement struct {
	ID           int         `gorm:"primaryKey" json:"-"`
	IngredientID int         `gorm:"not null" json:"-"`
	Ingredient   *Ingredient `gorm:"foreignKey:ingredient_id" json:"ingredient"`
	MovementType string      `gorm:"type:enum('in','out');not null;index" json:"movementType"`
	Quantity     int         `gorm:"not null;default:0" json:"quantity"`
	Description  string      `gorm:"size:100;null" json:"description"`
	CreatedAt    time.Time   `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	CreatedBy    int         `gorm:"column:created_by;not null;default:0" json:"-"`
}
