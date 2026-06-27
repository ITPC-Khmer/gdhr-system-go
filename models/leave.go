package models

import "time"

// Leave is a single leave request. status moves pending -> approved/rejected.
// On create the approval workflow (leave_approval rows) is seeded automatically;
// on approval the staff member's leave_year balance is updated. All *_id columns
// are plain indexed columns (no FK constraints), matching the rest of the schema.
type Leave struct {
	ID               int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	StaffID          string     `gorm:"size:100;not null;index" json:"staff_id" binding:"required"`
	RefNumber        string     `gorm:"size:700" json:"ref_number"`
	LeaveTypeID      int64      `gorm:"not null;index" json:"leave_type_id" binding:"required"`
	StartDate        Date       `gorm:"type:date;not null" json:"start_date"`
	EndDate          Date       `gorm:"type:date;not null" json:"end_date"`
	TotalDay         float64    `gorm:"default:0" json:"total_day"`
	Reason           string     `gorm:"size:1000" json:"reason"`
	Status           string     `gorm:"type:enum('pending','approved','rejected');default:'pending';index" json:"status"`
	ApprovedBy       *string    `gorm:"size:100;index" json:"approved_by"`
	ApprovedAt       *time.Time `json:"approved_at"`
	Attachment       string     `gorm:"type:text" json:"attachment"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	FileName         string     `gorm:"type:text" json:"file_name"`
	Seq              int        `gorm:"default:0" json:"seq"`
	RejectReason     string     `gorm:"size:1000" json:"reject_reason"`
	ApproveDocument  JSONRaw    `gorm:"type:json" json:"approve_document"`
	StatusDocument   string     `gorm:"type:enum('pending','completed');default:'pending'" json:"status_document"`
	RequiredDocument int8       `gorm:"default:0" json:"required_document"`
	ArrivedLv10      int8       `gorm:"default:0" json:"arrived_lv10"`
	GoAbroad         int        `gorm:"default:0" json:"go_abroad"`
	StaffsInfo       JSONRaw    `gorm:"type:json" json:"staffs_info"`
	Phone            string     `gorm:"size:700" json:"phone"`

	// LeaveTypeName is resolved at query time (leave_type_id -> type_name); it is
	// not a stored column.
	LeaveTypeName string `gorm:"-" json:"leave_type_name,omitempty"`
}

// TableName overrides the default table name.
func (Leave) TableName() string {
	return "leave"
}

// LeaveApproval is one step of a leave/mission approval workflow. Rows are seeded
// automatically when a leave request is created and flipped to approved/rejected
// as the request is processed.
type LeaveApproval struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	LType         string     `gorm:"column:l_type;type:enum('leave','mission');default:'leave';index" json:"l_type"`
	LeaveID       *int64     `gorm:"index;uniqueIndex:leave_id_2" json:"leave_id"`
	StaffID       *string    `gorm:"size:100;index;uniqueIndex:leave_id_2" json:"staff_id"`
	LLevel        *int       `gorm:"column:l_level;index" json:"l_level"`
	Status        string     `gorm:"type:enum('pending','approved','rejected');default:'pending';index" json:"status"`
	ShowType      *int       `gorm:"index" json:"show_type"`
	CreatedAt     *time.Time `gorm:"type:datetime" json:"created_at"`
	ApprovedAt    *time.Time `gorm:"type:datetime" json:"approved_at"`
	MaxLevel      int        `gorm:"default:1" json:"max_level"`
	ApproveLevel  *int       `gorm:"uniqueIndex:leave_id_2" json:"approve_level"`
	RoleName      string     `gorm:"size:100;default:'normal'" json:"role_name"`
	ModeratorType int8       `gorm:"default:0" json:"moderator_type"`
	InstituteID   *string    `gorm:"size:45" json:"institute_id"`
	N             int        `gorm:"default:0" json:"n"`
}

// TableName overrides the default table name.
func (LeaveApproval) TableName() string {
	return "leave_approval"
}

// LeaveRequesterFile is a comment/attachment thread entry attached to a leave
// request (chat messages or the request description), with files stored as JSON.
type LeaveRequesterFile struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	LeaveID       int64      `gorm:"not null;index" json:"leave_id" binding:"required"`
	Comment       string     `gorm:"type:text" json:"comment"`
	MessageType   string     `gorm:"type:enum('chat','description');not null;default:'description'" json:"message_type"`
	ParentID      *uint      `gorm:"index" json:"parent_id"`
	StaffID       *string    `gorm:"size:100;index" json:"staff_id"`
	TaskID        uint       `gorm:"not null;index" json:"task_id"`
	AddedBy       *string    `gorm:"size:100;index" json:"added_by"`
	LastUpdatedBy *string    `gorm:"size:100;index" json:"last_updated_by"`
	Status        string     `gorm:"type:enum('sent','read');not null;default:'sent'" json:"status"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	Files         JSONRaw    `gorm:"type:json" json:"files"`
}

// TableName overrides the default table name.
func (LeaveRequesterFile) TableName() string {
	return "leave_requester_files"
}

// LeaveYear is a per-staff, per-leave-type, per-year balance. It is upserted
// when a leave request is approved: total_day accumulates used days and
// leave_remaining tracks what is left against max_days.
type LeaveYear struct {
	ID             int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	StaffID        string  `gorm:"size:100;not null;uniqueIndex:uq_leave_year" json:"staff_id" binding:"required"`
	LeaveTypeID    int64   `gorm:"not null;index;uniqueIndex:uq_leave_year" json:"leave_type_id" binding:"required"`
	LYear          int     `gorm:"column:l_year;not null;uniqueIndex:uq_leave_year" json:"l_year" binding:"required"`
	TotalDay       float64 `gorm:"default:0" json:"total_day"`
	IsReset        bool    `gorm:"default:0" json:"is_reset"`
	MaxDays        int     `gorm:"default:0" json:"max_days"`
	LeaveRemaining float64 `gorm:"default:0" json:"leave_remaining"`
}

// TableName overrides the default table name.
func (LeaveYear) TableName() string {
	return "leave_year"
}