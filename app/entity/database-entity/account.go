package databaseentity

import "time"

// table accounts model
type Account struct {
	ID        int       `gorm:"primaryKey" json:"-"`
	Name      string    `gorm:"unique;size:100;not null" json:"name"`
	Type      string    `gorm:"type:enum('asset', 'liability', 'equity', 'revenue', 'expense');not null;index" json:"type"`
	Balance   float32   `gorm:"not null;type:decimal(10,2);default:0" json:"balance"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	CreatedBy int       `gorm:"column:created_by;not null;default:0" json:"-"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"updatedAt"`
	UpdatedBy int       `gorm:"column:updated_by;not null;default:0" json:"-"`
}
