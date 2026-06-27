package models

import "time"

// Position mirrors a GDHR public-positions record. position_id is the natural
// primary key (provided by the source, not auto-incremented).
type Position struct {
	PositionID        int       `gorm:"primaryKey;autoIncrement:false" json:"position_id"`
	PositionName      string    `gorm:"size:255" json:"position_name"`
	Description       string    `gorm:"size:512" json:"description"`
	PositionNameShort string    `gorm:"size:128" json:"position_name_short"`
	PositionBaseID    int       `gorm:"index" json:"position_base_id"`
	RankBaseID        int       `gorm:"index" json:"rank_base_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TableName overrides the default table name.
func (Position) TableName() string {
	return "positions"
}
