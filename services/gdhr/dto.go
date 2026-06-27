package gdhr

import (
	"bytes"
	"encoding/json"
)

// FlexBool decodes a boolean that the GDHR API expresses in several shapes:
//   - a Node Buffer object: {"type":"Buffer","data":[1]}  (ranks.active)
//   - a string: "1" / "0"                                  (institutes.active)
//   - a plain number or JSON bool
type FlexBool bool

// UnmarshalJSON implements json.Unmarshaler for the shapes above.
func (b *FlexBool) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		*b = false
		return nil
	}

	switch data[0] {
	case '{': // Buffer object: {"type":"Buffer","data":[1]}
		var buf struct {
			Data []int `json:"data"`
		}
		if err := json.Unmarshal(data, &buf); err != nil {
			return err
		}
		*b = FlexBool(len(buf.Data) > 0 && buf.Data[0] != 0)
	case '"': // string "1" / "0" / "true"
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*b = FlexBool(s == "1" || s == "true")
	case 't', 'f': // real bool
		var v bool
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		*b = FlexBool(v)
	default: // number
		var n float64
		if err := json.Unmarshal(data, &n); err != nil {
			return err
		}
		*b = FlexBool(n != 0)
	}
	return nil
}

// ---- institutes: { page, totalPages, totalItems, items[] } ----

type InstituteDTO struct {
	ID            string   `json:"id"`
	OldID         string   `json:"old_id"`
	ParentID      string   `json:"parent_id"`
	LevelType     int      `json:"level_type"`
	OldParentID   int      `json:"old_parent_id"`
	Name          string   `json:"name"`
	NameShort     string   `json:"name_short"`
	Active        FlexBool `json:"active"`
	SourceTable   string   `json:"source_table"`
	SourceTableID int      `json:"source_table_id"`
}

type InstitutesResponse struct {
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
	TotalPages int            `json:"totalPages"`
	TotalItems int            `json:"totalItems"`
	Items      []InstituteDTO `json:"items"`
}

// ---- staffs: { total, page, pageSize, data[] } ----

type StaffDTO struct {
	UID                          string `json:"uid"`
	StaffTypeID                  int    `json:"staff_type_id"`
	StaffNumber                  string `json:"staff_number"`
	SurnameKh                    string `json:"surname_kh"`
	NameKh                       string `json:"name_kh"`
	SurnameEn                    string `json:"surname_en"`
	NameEn                       string `json:"name_en"`
	PlaceOfBirth                 string `json:"place_of_birth"`
	Nationality                  string `json:"nationality"`
	Address                      string `json:"address"`
	PhotoPath                    string `json:"photo_path"`
	Gender                       int    `json:"gender"`
	Phone                        string `json:"phone"`
	Email                        string `json:"email"`
	CityzenCardNumber            string `json:"cityzen_card_number"`
	RankID                       int    `json:"rank_id"`
	RankNameShort                string `json:"rank_name_short"`
	PositionID                   int    `json:"position_id"`
	PositionNameShort            string `json:"position_name_short"`
	OtherPositionID              int    `json:"other_position_id"`
	GeneralCommissariatID        int    `json:"general_commissariat_id"`
	GeneralCommissariatNameShort string `json:"general_commissariat_name_short"`
	DepartmentID                 int    `json:"department_id"`
	DepartmentNameShort          string `json:"department_name_short"`
	OfficeID                     int    `json:"office_id"`
	OfficeNameShort              string `json:"office_name_short"`
	SectorID                     int    `json:"sector_id"`
	SectorNameShort              string `json:"sector_name_short"`
	InstituteID                  string `json:"institute_id"`
	StatusID                     int    `json:"status_id"`
	StatusName                   string `json:"status_name"`
}

type StaffsResponse struct {
	Total    int        `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
	Data     []StaffDTO `json:"data"`
}

// ---- ranks: { data: { items[] } } ----

type RankDTO struct {
	RankID          int      `json:"rank_id"`
	RankName        string   `json:"rank_name"`
	RankNameShort   string   `json:"rank_name_short"`
	PositionBaseID  int      `json:"position_base_id"`
	RankOrder       int      `json:"rank_order"`
	PromotePeriod   int      `json:"promote_peroid"` // misspelled in source
	RankNameEn      string   `json:"rank_name_en"`
	RankNameShortEn string   `json:"rank_name_short_en"`
	Active          FlexBool `json:"active"`
}

type RanksResponse struct {
	Data struct {
		Items []RankDTO `json:"items"`
	} `json:"data"`
}

// ---- positions: { data: { items[] } } ----

type PositionDTO struct {
	PositionID        int    `json:"position_id"`
	PositionName      string `json:"position_name"`
	Description       string `json:"description"`
	PositionNameShort string `json:"position_name_short"`
	PositionBaseID    int    `json:"position_base_id"`
	RankBaseID        int    `json:"rank_base_id"`
}

type PositionsResponse struct {
	Data struct {
		Items []PositionDTO `json:"items"`
	} `json:"data"`
}
