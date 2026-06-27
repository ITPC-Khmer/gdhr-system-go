package models

import "time"

// Staff mirrors a GDHR public-staffs record. All *_id fields are plain indexed
// columns, NOT foreign keys.
type Staff struct {
	UID                          string    `gorm:"primaryKey;size:100" json:"uid"`
	StaffTypeID                  int       `gorm:"index" json:"staff_type_id"`
	StaffNumber                  string    `gorm:"index;size:64" json:"staff_number"`
	SurnameKh                    string    `gorm:"size:255" json:"surname_kh"`
	NameKh                       string    `gorm:"size:255" json:"name_kh"`
	SurnameEn                    string    `gorm:"size:255" json:"surname_en"`
	NameEn                       string    `gorm:"size:255" json:"name_en"`
	PlaceOfBirth                 string    `gorm:"size:255" json:"place_of_birth"`
	Nationality                  string    `gorm:"size:191" json:"nationality"`
	Address                      string    `gorm:"size:512" json:"address"`
	PhotoPath                    string    `gorm:"size:512" json:"photo_path"`
	Gender                       int       `json:"gender"`
	Phone                        string    `gorm:"size:191" json:"phone"`
	Email                        string    `gorm:"size:191" json:"email"`
	CityzenCardNumber            string    `gorm:"size:100" json:"cityzen_card_number"`
	RankID                       int       `gorm:"index" json:"rank_id"`
	RankNameShort                string    `gorm:"size:191" json:"rank_name_short"`
	PositionID                   int       `gorm:"index" json:"position_id"`
	PositionNameShort            string    `gorm:"size:191" json:"position_name_short"`
	OtherPositionID              int       `gorm:"index" json:"other_position_id"`
	GeneralCommissariatID        int       `gorm:"index" json:"general_commissariat_id"`
	GeneralCommissariatNameShort string    `gorm:"size:255" json:"general_commissariat_name_short"`
	DepartmentID                 int       `gorm:"index" json:"department_id"`
	DepartmentNameShort          string    `gorm:"size:255" json:"department_name_short"`
	OfficeID                     int       `gorm:"index" json:"office_id"`
	OfficeNameShort              string    `gorm:"size:255" json:"office_name_short"`
	SectorID                     int       `gorm:"index" json:"sector_id"`
	SectorNameShort              string    `gorm:"size:255" json:"sector_name_short"`
	InstituteID                  string    `gorm:"index;size:100" json:"institute_id"`
	StatusID                     int       `gorm:"index" json:"status_id"`
	StatusName                   string    `gorm:"size:191" json:"status_name"`
	CreatedAt                    time.Time `json:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at"`

	// Denormalized institute chain (self -> parent -> ... -> root), recomputed
	// from institute_id on insert/update and during sync. Stored as JSON.
	InstituteIDs       StringSlice `gorm:"type:json" json:"institute_ids"`
	InstituteHierarchy StringSlice `gorm:"type:json" json:"institute_hierarchy"`

	// InstituteName is resolved at query time (institute_id -> institute's name);
	// it is not a stored column.
	InstituteName string `gorm:"-" json:"institute_name,omitempty"`
}

// TableName overrides the default table name.
func (Staff) TableName() string {
	return "staffs"
}
