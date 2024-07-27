package databaseentity

import "time"

// table units model
type Unit struct {
	ID        int       `gorm:"primaryKey" json:"-"`
	Name      string    `gorm:"unique;size:30;not null" json:"name"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	CreatedBy int       `gorm:"column:created_by;not null;default:0" json:"-"`
}
