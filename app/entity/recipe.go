package entity

import "time"

// table recipes model
type Recipe struct {
	ID                  int       `gorm:"primaryKey" json:"-"`
	Serial              string    `gorm:"unique;size:11;not null" json:"serial"`
	Name                string    `gorm:"size:100;not null" json:"name"`
	Quantity            int       `gorm:"not null;default:0" json:"quantity"`
	Description         string    `gorm:"type:text;null" json:"description"`
	LaborDescription    string    `gorm:"type:text;null" json:"laborDescription"`
	OverheadDescription string    `gorm:"type:text;null" json:"overheadDescription"`
	RawMaterialCosts    float32   `gorm:"not null;type:decimal(10,2);default:0" json:"rawMaterialCosts"`
	LaborCosts          float32   `gorm:"not null;type:decimal(10,2);default:0" json:"laborCosts"`
	OverheadCosts       float32   `gorm:"not null;type:decimal(10,2);default:0" json:"overheadCosts"`
	ExpectedProfit      int       `gorm:"not null;default:0" json:"expectedProfit"`
	SellingPrice        float32   `gorm:"-" json:"sellingPrice"`
	CreatedAt           time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	CreatedBy           string    `gorm:"size:30;not null;default:'SYSTEM'" json:"createdBy"`
	UpdatedAt           time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"updatedAt"`
	UpdatedBy           string    `gorm:"size:30;not null;default:'SYSTEM'" json:"updatedBy"`
}

func (e *Recipe) SetSellingPrice() {
	hpp := e.RawMaterialCosts + e.LaborCosts + e.OverheadCosts
	e.SellingPrice = hpp + (float32(e.ExpectedProfit) / float32(100) * hpp)
}
