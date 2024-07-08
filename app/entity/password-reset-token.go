package entity

import "time"

type PasswordResetToken struct {
	ID        int       `gorm:"primaryKey" json:"-"`
	Email     string    `gorm:"unique;size:255;not null" json:"email"`
	Token     string    `gorm:"size:255;null" json:"token"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"updatedAt"`
}
