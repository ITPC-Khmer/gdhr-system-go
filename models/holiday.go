package models

import (
	"crypto/rand"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Holiday represents a public/company holiday. holiday_id is a UUID (char 36)
// generated on insert. date is the canonical day; from_date/to_date describe an
// optional multi-day span. created_by_id/updated_by_id reference users.id.
type Holiday struct {
	HolidayID   string    `gorm:"primaryKey;type:char(36)" json:"holiday_id"`
	Name        string    `gorm:"size:255;not null" json:"name" binding:"required"`
	Date        Date      `gorm:"type:date;not null;index:idx_holidays_date" json:"date"`
	IsRecurring bool      `gorm:"not null;default:0" json:"is_recurring"`
	FromDate    *Date     `gorm:"type:date" json:"from_date"`
	ToDate      *Date     `gorm:"type:date" json:"to_date"`
	CreatedByID *uint     `gorm:"index" json:"created_by_id"`
	UpdatedByID *uint     `gorm:"index" json:"updated_by_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName overrides the default table name.
func (Holiday) TableName() string {
	return "holidays"
}

// BeforeCreate assigns a UUID primary key when one wasn't supplied.
func (h *Holiday) BeforeCreate(*gorm.DB) error {
	if h.HolidayID == "" {
		id, err := newUUID()
		if err != nil {
			return err
		}
		h.HolidayID = id
	}
	return nil
}

// newUUID returns a random RFC-4122 v4 UUID string.
func newUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
