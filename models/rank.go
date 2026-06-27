package models

import "time"

// Rank mirrors a GDHR public-ranks record. rank_id is the natural primary key
// (provided by the source, not auto-incremented).
type Rank struct {
	RankID          int       `gorm:"primaryKey;autoIncrement:false" json:"rank_id"`
	RankName        string    `gorm:"size:255" json:"rank_name"`
	RankNameShort   string    `gorm:"size:128" json:"rank_name_short"`
	PositionBaseID  int       `gorm:"index" json:"position_base_id"`
	RankOrder       int       `json:"rank_order"`
	PromotePeriod   int       `json:"promote_period"`
	RankNameEn      string    `gorm:"size:255" json:"rank_name_en"`
	RankNameShortEn string    `gorm:"size:128" json:"rank_name_short_en"`
	Active          bool      `json:"active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName overrides the default table name.
func (Rank) TableName() string {
	return "ranks"
}
