package models

import "time"

// LeaveType is a category of leave (annual, short, maternity, …). max_days is
// the per-year allowance used to seed/refresh a staff member's leave_year
// balance; is_reset marks types whose balance resets each year.
type LeaveType struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	LeaveKey    string    `gorm:"size:255" json:"leave_key"`
	TypeName    string    `gorm:"size:255;not null;uniqueIndex" json:"type_name" binding:"required"`
	TypeNameS   string    `gorm:"size:255;not null" json:"type_name_s" binding:"required"`
	Description string    `gorm:"size:1000" json:"description"`
	IsReset     bool      `gorm:"default:0" json:"is_reset"`
	MaxDays     *int      `json:"max_days"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName overrides the default table name.
func (LeaveType) TableName() string {
	return "leave_type"
}