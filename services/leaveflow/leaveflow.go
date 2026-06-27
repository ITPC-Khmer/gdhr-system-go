// Package leaveflow resolves WHO must approve a leave request by walking the
// requester's institute chain (self -> ... -> root) and matching managers from
// staff_institute_roles. See docs/LEAVE-SYSTEM-SCHEMA.txt (PROCESS FLOW).
//
// Authority is encoded by institutes.level_type: SMALLER level_type = HIGHER
// authority. The requester's institute chain is therefore an ordered list of
// tiers from the requester's own level (largest number) up to the root (10).
package leaveflow

import (
	"backend/models"
	"backend/services/hierarchy"

	"gorm.io/gorm"
)

// Approver is a resolved approval step: which staff, at which institute/tier,
// holding which management role.
type Approver struct {
	StaffID     string
	InstituteID string
	LevelType   int
	Role        string // "moderator" | "approval"
}

// tier is one institute in the requester's chain with its authority level.
type tier struct {
	InstituteID string
	Level       int
}

// requesterChain resolves the requesting staff's institute tiers, ordered from
// their own level toward higher authority (self -> ... -> root). It returns the
// requester's own level_type, the tiers, and ok=false when the staff or their
// institutes can't be resolved.
func requesterChain(db *gorm.DB, staffID string) (ownLevel int, tiers []tier, ok bool) {
	var staff models.Staff
	if err := db.Select("uid", "institute_id", "institute_ids").
		Where("uid = ?", staffID).First(&staff).Error; err != nil {
		return 0, nil, false
	}

	ids := []string(staff.InstituteIDs)
	if len(ids) == 0 {
		// Fall back to resolving the chain live from parent_id links.
		ids, _ = hierarchy.Chain(db, staff.InstituteID)
	}
	if len(ids) == 0 {
		return 0, nil, false
	}

	var insts []models.Institute
	db.Select("id", "level_type").Where("id IN ?", ids).Find(&insts)
	levelByID := make(map[string]int, len(insts))
	for _, in := range insts {
		levelByID[in.ID] = in.LevelType
	}

	tiers = make([]tier, 0, len(ids))
	for _, id := range ids { // ids are already self -> parent -> ... -> root
		if lv, exists := levelByID[id]; exists {
			tiers = append(tiers, tier{InstituteID: id, Level: lv})
		}
	}
	if len(tiers) == 0 {
		return 0, nil, false
	}
	return tiers[0].Level, tiers, true
}

// findManager returns the staff managing instituteID with the given role, or nil.
func findManager(db *gorm.DB, instituteID, role string, level int) *Approver {
	var r models.StaffInstituteRole
	if err := db.Where("institute_id = ? AND role = ?", instituteID, role).
		First(&r).Error; err != nil {
		return nil
	}
	return &Approver{StaffID: r.StaffID, InstituteID: instituteID, LevelType: level, Role: role}
}

// walk scans tiers at or above fromLevel authority (Level <= fromLevel), in
// chain order, returning the first matching manager. When includeModerator is
// true a tier's moderator is preferred over its approval.
func walk(db *gorm.DB, tiers []tier, fromLevel int, includeModerator bool) *Approver {
	for _, t := range tiers {
		if t.Level > fromLevel {
			continue // lower authority than the start tier — skip
		}
		if includeModerator {
			if a := findManager(db, t.InstituteID, "moderator", t.Level); a != nil {
				return a
			}
		}
		if a := findManager(db, t.InstituteID, "approval", t.Level); a != nil {
			return a
		}
	}
	return nil
}

// RequesterLevel returns the requesting staff's own level_type.
func RequesterLevel(db *gorm.DB, staffID string) (int, bool) {
	own, _, ok := requesterChain(db, staffID)
	return own, ok
}

// FirstApprover resolves the initial approver for a new request: starting at the
// requester's own tier and walking toward higher authority, moderator preferred
// over approval, first match wins. Returns nil when no manager is configured.
func FirstApprover(db *gorm.DB, requesterStaffID string) *Approver {
	own, tiers, ok := requesterChain(db, requesterStaffID)
	if !ok {
		return nil
	}
	return walk(db, tiers, own, true)
}

// NextApproval resolves the next 'approval'-role approver, starting at fromLevel
// (the current step's tier) and walking toward higher authority. Returns nil
// when no further approval-role manager exists up to the root.
func NextApproval(db *gorm.DB, requesterStaffID string, fromLevel int) *Approver {
	_, tiers, ok := requesterChain(db, requesterStaffID)
	if !ok {
		return nil
	}
	return walk(db, tiers, fromLevel, false)
}
