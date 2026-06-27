package models

import "time"

// User is the application user / account record.
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:120;not null" json:"name"`
	Email     string    `gorm:"size:160;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"` // never serialized
	Role      string    `gorm:"size:40;not null;default:user" json:"role"`
	Active    bool      `gorm:"not null;default:true" json:"active"`
	// StaffID optionally links this account to a staff member (staffs.uid). Used
	// to authorize per-approver leave actions (the caller must be the assigned
	// approver). Empty = not linked.
	StaffID   string    `gorm:"size:100;index" json:"staff_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the default table name.
func (User) TableName() string {
	return "users"
}
