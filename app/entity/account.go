package entity

import "time"

// table accounts model
type Account struct {
	ID        int       `gorm:"primaryKey" json:"-"`
	Name      string    `gorm:"unique;size:100;not null" json:"name"`
	Type      string    `gorm:"type:enum('asset', 'liability', 'equity', 'revenue', 'expense');not null;index" json:"type"`
	Balance   float32   `gorm:"not null;type:decimal(10,2);default:0" json:"balance"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	CreatedBy string    `gorm:"size:30;not null;default:'SYSTEM'" json:"createdBy"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"updatedAt"`
	UpdatedBy string    `gorm:"size:30;not null;default:'SYSTEM'" json:"updatedBy"`
}
