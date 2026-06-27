package models

import "time"

// StaffInstituteRole maps a staff member to an institute they manage, with a
// role. A staff can manage many institutes and hold several roles; this is the
// many-to-many mapping the leave-approval routing reads (match institute by
// level_type, then pick a moderator/approval manager).
//
// Business rule (enforced in application logic, not by a DB constraint): within
// one institute there is at most ONE moderator and at most ONE approval; admin
// may be assigned to many. staff_id / institute_id are plain indexed columns
// (no FK), matching the rest of this codebase.
type StaffInstituteRole struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	StaffID     string    `gorm:"size:100;not null;index;uniqueIndex:uq_staff_institute_role" json:"staff_id" binding:"required"`
	InstituteID string    `gorm:"size:100;not null;index;uniqueIndex:uq_staff_institute_role" json:"institute_id" binding:"required"`
	Role        string    `gorm:"type:enum('admin','moderator','approval');not null;index;uniqueIndex:uq_staff_institute_role" json:"role" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// InstituteName is resolved at query time (institute_id -> institute's name);
	// it is not a stored column.
	InstituteName string `gorm:"-" json:"institute_name,omitempty"`
}

// TableName overrides the default table name.
func (StaffInstituteRole) TableName() string {
	return "staff_institute_roles"
}
