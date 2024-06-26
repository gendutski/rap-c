package entity

import "time"

// table users model
type User struct {
	ID                 int       `gorm:"primaryKey" json:"-"`
	Username           string    `gorm:"unique;size:30;not null" json:"userName"`
	FullName           string    `gorm:"size:100;not null" json:"fullName"`
	Email              string    `gorm:"unique;size:100;not null" json:"email"`
	Password           string    `gorm:"size:255;not null" json:"-"`
	PasswordMustChange bool      `gorm:"not null;default:0" json:"passwordMustChange"`
	Disabled           bool      `gorm:"not null;default:0" json:"disabled"`
	IsGuest            bool      `gorm:"not null;default:0" json:"isGuest"`
	CreatedAt          time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	CreatedBy          string    `gorm:"size:30;not null;default:'SYSTEM'" json:"createdBy"`
	UpdatedAt          time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"updatedAt"`
	UpdatedBy          string    `gorm:"size:30;not null;default:'SYSTEM'" json:"updatedBy"`
}
