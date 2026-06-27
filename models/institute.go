package models

import "time"

// Institute mirrors a GDHR public-institutes record. parent_id is NOT a foreign
// key — it's a plain indexed column.
type Institute struct {
	ID            string    `gorm:"primaryKey;size:100" json:"id"`
	OldID         string    `gorm:"size:32" json:"old_id"`
	ParentID      string    `gorm:"index;size:100" json:"parent_id"`
	LevelType     int       `gorm:"index" json:"level_type"`
	OldParentID   int       `json:"old_parent_id"`
	Name          string    `gorm:"size:255" json:"name"`
	NameShort     string    `gorm:"size:255" json:"name_short"`
	Active        bool      `json:"active"`
	SourceTable   string    `gorm:"size:64" json:"source_table"`
	SourceTableID int       `json:"source_table_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// ParentName is resolved at query time (parent_id -> parent's name); it is
	// not a stored column.
	ParentName string `gorm:"-" json:"parent_name,omitempty"`
}

// TableName overrides the default table name.
func (Institute) TableName() string {
	return "institutes"
}
