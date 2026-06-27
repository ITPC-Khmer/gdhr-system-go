package models

// LeaveRole encodes the approval/duration rules for a leave_type. For a given
// requested duration (min_duration..max_duration) and staff_type, approve_level
// drives how many approval steps the request must pass through. leave_type_id is
// a plain indexed column (no FK constraint, matching the rest of this codebase).
type LeaveRole struct {
	ID              int    `gorm:"primaryKey;autoIncrement" json:"id"`
	LeaveTypeID     *int64 `gorm:"index" json:"leave_type_id"`
	LeaveType       string `gorm:"column:leave_type;size:50;not null" json:"leave_type" binding:"required"`
	Name            string `gorm:"size:255;not null" json:"name" binding:"required"`
	MinDuration     int    `gorm:"not null" json:"min_duration"`
	LimitDays       int    `gorm:"not null" json:"limit_days"`
	MinDurationShow int    `gorm:"not null;default:1" json:"min_duration_show"`
	MaxDuration     int    `gorm:"not null" json:"max_duration"`
	ApproveLevel    int    `gorm:"not null" json:"approve_level"`
	StaffType       string `gorm:"size:50;not null" json:"staff_type" binding:"required"`
}

// TableName overrides the default table name.
func (LeaveRole) TableName() string {
	return "leave_roles"
}