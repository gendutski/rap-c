package databaseentity

import "time"

// table transactions model
type Transaction struct {
	ID          int       `gorm:"primaryKey" json:"-"`
	AccountID   int       `gorm:"not null" json:"-"`
	Account     *Account  `gorm:"foreignKey:account_id" json:"account"`
	Type        string    `gorm:"type:enum('debit', 'credit');not null;index" json:"type"`
	Amount      float32   `gorm:"not null;type:decimal(10,2);default:0" json:"amount"`
	Description string    `gorm:"type:text;null" json:"description"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	CreatedBy   int       `gorm:"column:created_by;not null;default:0" json:"-"`
}
