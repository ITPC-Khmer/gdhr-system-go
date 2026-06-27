package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"backend/database"
	"backend/models"
	"backend/services/leaveflow"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CRUD + workflow for the leave request family:
//   - leave                 the request itself (create seeds the approval workflow)
//   - leave_approval        per-step approval rows (auto-seeded on request)
//   - leave_requester_files comment/attachment thread entries
//   - leave_year            per-staff/per-type/per-year balance (upserted on approve)
//
// Two automations live here:
//   1. Creating a leave request auto-seeds its leave_approval workflow ("when
//      request live => auto").
//   2. Approving a leave request upserts the requester's leave_year balance.

// ---- leave ----

var leaveSearch = []string{"staff_id", "ref_number", "reason", "phone"}
var leaveFilter = []string{"staff_id", "leave_type_id", "status", "approved_by", "status_document"}

// enrichLeaves resolves each row's leave_type_id to the leave type's name in a
// single query (leave_type_id is a plain column, not a GORM relation).
func enrichLeaves(rows []models.Leave) {
	ids := make([]int64, 0, len(rows))
	seen := map[int64]bool{}
	for _, r := range rows {
		if r.LeaveTypeID != 0 && !seen[r.LeaveTypeID] {
			seen[r.LeaveTypeID] = true
			ids = append(ids, r.LeaveTypeID)
		}
	}
	if len(ids) == 0 {
		return
	}
	var types []models.LeaveType
	database.DB.Select("id", "type_name").Where("id IN ?", ids).Find(&types)
	names := make(map[int64]string, len(types))
	for _, t := range types {
		names[t.ID] = t.TypeName
	}
	for i := range rows {
		rows[i].LeaveTypeName = names[rows[i].LeaveTypeID]
	}
}

func ListLeaves(c *gin.Context) {
	page, limit, offset := paginate(c)

	q := applyListFilters(c, database.DB.Model(&models.Leave{}), leaveSearch, leaveFilter)

	var total int64
	q.Count(&total)

	rows := make([]models.Leave, 0, limit)
	if err := q.Order("id DESC").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch records"})
		return
	}
	enrichLeaves(rows)

	c.JSON(http.StatusOK, gin.H{"data": rows, "total": total, "page": page, "limit": limit})
}

func GetLeave(c *gin.Context) { getResource[models.Leave](c, "id") }

// LeaveStats returns the leave-request counts by status (plus total) in a single
// grouped query — used by the list view's summary cards.
func LeaveStats(c *gin.Context) {
	type row struct {
		Status string
		Count  int64
	}
	var rows []row
	if err := database.DB.Model(&models.Leave{}).
		Select("status, COUNT(*) AS count").Group("status").Scan(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch stats"})
		return
	}

	out := gin.H{"total": int64(0), "pending": int64(0), "approved": int64(0), "rejected": int64(0)}
	var total int64
	for _, r := range rows {
		out[r.Status] = r.Count
		total += r.Count
	}
	out["total"] = total
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// CreateLeave records the request and, in the same transaction, auto-seeds its
// approval workflow so the request goes "live" ready for processing.
func CreateLeave(c *gin.Context) {
	var lv models.Leave
	if err := c.ShouldBindJSON(&lv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	lv.ID = 0 // never trust a client-supplied id
	if lv.Status == "" {
		lv.Status = "pending"
	}
	// Derive total_day (inclusive of both endpoints) when the client omits it.
	if lv.TotalDay == 0 && !lv.StartDate.Time.IsZero() && !lv.EndDate.Time.IsZero() {
		if d := lv.EndDate.Time.Sub(lv.StartDate.Time).Hours()/24 + 1; d > 0 {
			lv.TotalDay = d
		}
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&lv).Error; err != nil {
			return err
		}
		return seedLeaveApprovals(tx, &lv)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create leave request"})
		return
	}

	enrichLeaves([]models.Leave{lv})
	c.JSON(http.StatusCreated, gin.H{"data": lv})
}

// seedLeaveApprovals creates the initial (pending) approval step for a freshly
// created request by resolving the first approver via the routing walk (see
// services/leaveflow). When no manager is configured in the requester's chain
// the request is left with no approval rows for an admin to handle manually.
func seedLeaveApprovals(tx *gorm.DB, lv *models.Leave) error {
	approver := leaveflow.FirstApprover(tx, lv.StaffID)
	if approver == nil {
		return nil
	}
	return insertApproval(tx, lv, approver)
}

// insertApproval writes one pending leave_approval task for the resolved
// approver. l_level / approve_level carry the approver's institute level_type
// (the tier); role_name carries the management role (moderator/approval).
func insertApproval(tx *gorm.DB, lv *models.Leave, approver *leaveflow.Approver) error {
	now := time.Now()
	level := approver.LevelType
	staffID := approver.StaffID
	instituteID := approver.InstituteID
	approval := models.LeaveApproval{
		LType:        "leave",
		LeaveID:      &lv.ID,
		StaffID:      &staffID,
		LLevel:       &level,
		Status:       "pending",
		CreatedAt:    &now,
		MaxLevel:     1,
		ApproveLevel: &level,
		RoleName:     approver.Role,
		InstituteID:  &instituteID,
	}
	return tx.Create(&approval).Error
}

// upsertLeaveYear adds the approved request's days onto the staff member's
// balance for the leave's year, creating the row on first use. max_days and
// is_reset come from the leave_type; leave_remaining = max_days - total used.
func upsertLeaveYear(tx *gorm.DB, lv *models.Leave) error {
	year := time.Now().Year()
	if !lv.StartDate.Time.IsZero() {
		year = lv.StartDate.Time.Year()
	}

	maxDays := 0
	isReset := false
	var lt models.LeaveType
	if err := tx.Where("id = ?", lv.LeaveTypeID).First(&lt).Error; err == nil {
		if lt.MaxDays != nil {
			maxDays = *lt.MaxDays
		}
		isReset = lt.IsReset
	}

	ly := models.LeaveYear{
		StaffID:        lv.StaffID,
		LeaveTypeID:    lv.LeaveTypeID,
		LYear:          year,
		TotalDay:       lv.TotalDay,
		IsReset:        isReset,
		MaxDays:        maxDays,
		LeaveRemaining: float64(maxDays) - lv.TotalDay,
	}
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "staff_id"}, {Name: "leave_type_id"}, {Name: "l_year"}},
		DoUpdates: clause.Assignments(map[string]any{
			"total_day":       gorm.Expr("leave_year.total_day + ?", lv.TotalDay),
			"max_days":        maxDays,
			"is_reset":        isReset,
			"leave_remaining": gorm.Expr("? - (leave_year.total_day + ?)", maxDays, lv.TotalDay),
		}),
	}).Create(&ly).Error
}

// ---- admin break-glass override ----
//
// Admins are otherwise view-only on leave. These endpoints exist ONLY to
// force-resolve a STUCK request (e.g. no approver could be routed, or an
// assigned approver has no usable account). They are admin-only, force the
// final state directly, and write an audit stamp into approve_document.

// overrideStamp merges the audit record for a break-glass action into the
// existing approve_document (preserving any other keys) and returns the updated
// JSON plus the acting admin's identity (email, falling back to "admin-override").
func overrideStamp(c *gin.Context, existing models.JSONRaw, action, note string, at time.Time) (models.JSONRaw, string) {
	by, _ := c.Get("email")
	actor, _ := by.(string)
	if actor == "" {
		actor = "admin-override"
	}
	// Start from the existing document when it's a JSON object; otherwise begin
	// fresh (a non-object value can't carry a merged key).
	doc := map[string]any{}
	if len(existing) > 0 {
		if err := json.Unmarshal(existing, &doc); err != nil {
			doc = map[string]any{}
		}
	}
	doc["override"] = map[string]any{
		"by":     actor,
		"action": action,
		"note":   note,
		"at":     at.Format(time.RFC3339),
	}
	b, _ := json.Marshal(doc)
	return models.JSONRaw(b), actor
}

// OverrideApproveLeave force-approves a stuck request: marks any pending steps
// approved, finalizes the leave, and upserts the year balance — all atomically.
func OverrideApproveLeave(c *gin.Context) {
	var lv models.Leave
	if err := database.DB.Where("id = ?", c.Param("id")).First(&lv).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "leave not found"})
		return
	}

	var body struct {
		Note string `json:"note"`
	}
	_ = c.ShouldBindJSON(&body)

	now := time.Now()
	doc, actor := overrideStamp(c, lv.ApproveDocument, "approve", body.Note, now)
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		lv.Status = "approved"
		lv.ApprovedAt = &now
		lv.ApprovedBy = &actor
		lv.ApproveDocument = doc
		if err := tx.Save(&lv).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.LeaveApproval{}).
			Where("leave_id = ? AND status = ?", lv.ID, "pending").
			Updates(map[string]any{"status": "approved", "approved_at": now}).Error; err != nil {
			return err
		}
		return upsertLeaveYear(tx, &lv)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to override-approve leave"})
		return
	}

	enrichLeaves([]models.Leave{lv})
	c.JSON(http.StatusOK, gin.H{"data": lv})
}

// OverrideRejectLeave force-rejects a stuck request: marks any pending steps
// rejected and finalizes the leave with the given reason.
func OverrideRejectLeave(c *gin.Context) {
	var lv models.Leave
	if err := database.DB.Where("id = ?", c.Param("id")).First(&lv).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "leave not found"})
		return
	}

	var body struct {
		RejectReason string `json:"reject_reason"`
	}
	_ = c.ShouldBindJSON(&body)

	now := time.Now()
	doc, actor := overrideStamp(c, lv.ApproveDocument, "reject", body.RejectReason, now)
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		lv.Status = "rejected"
		lv.RejectReason = body.RejectReason
		lv.ApprovedAt = &now
		lv.ApprovedBy = &actor
		lv.ApproveDocument = doc
		if err := tx.Save(&lv).Error; err != nil {
			return err
		}
		return tx.Model(&models.LeaveApproval{}).
			Where("leave_id = ? AND status = ?", lv.ID, "pending").
			Updates(map[string]any{"status": "rejected", "approved_at": now}).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to override-reject leave"})
		return
	}

	enrichLeaves([]models.Leave{lv})
	c.JSON(http.StatusOK, gin.H{"data": lv})
}

// ---- leave_approval ----

var leaveApprovalSearch = []string{"staff_id", "role_name", "institute_id"}
var leaveApprovalFilter = []string{"l_type", "leave_id", "staff_id", "status", "l_level", "approve_level"}

func ListLeaveApprovals(c *gin.Context) {
	listResource[models.LeaveApproval](c, leaveApprovalSearch, leaveApprovalFilter, "id DESC")
}
func GetLeaveApproval(c *gin.Context) { getResource[models.LeaveApproval](c, "id") }

// callerStaffID returns the staff_id linked to the authenticated user account,
// or "" when the user can't be resolved or isn't linked to a staff member.
func callerStaffID(c *gin.Context) string {
	uid := authUserID(c)
	if uid == nil {
		return ""
	}
	var u models.User
	if err := database.DB.Select("staff_id").Where("id = ?", *uid).First(&u).Error; err != nil {
		return ""
	}
	return u.StaffID
}

// authorizeApprover ensures the caller is the staff member assigned to this
// approval task. Admins have VIEW-ONLY access to leave and may NOT act. Writes a
// 403 and returns false otherwise.
func authorizeApprover(c *gin.Context, ap *models.LeaveApproval) bool {
	if role, _ := c.Get("role"); role == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"message": "admins have view-only access to leave; only the assigned approver can act"})
		return false
	}
	sid := callerStaffID(c)
	if sid == "" || ap.StaffID == nil || sid != *ap.StaffID {
		c.JSON(http.StatusForbidden, gin.H{"message": "you are not the assigned approver for this task"})
		return false
	}
	return true
}

// loadPendingStep loads a pending leave_approval task plus its parent leave, or
// writes the appropriate error response and returns ok=false.
func loadPendingStep(c *gin.Context) (ap models.LeaveApproval, lv models.Leave, ok bool) {
	if err := database.DB.Where("id = ?", c.Param("id")).First(&ap).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "approval task not found"})
		return ap, lv, false
	}
	if ap.Status != "pending" {
		c.JSON(http.StatusConflict, gin.H{"message": "this task is already " + ap.Status})
		return ap, lv, false
	}
	if ap.LeaveID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "task is not linked to a leave"})
		return ap, lv, false
	}
	if err := database.DB.Where("id = ?", *ap.LeaveID).First(&lv).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "leave not found"})
		return ap, lv, false
	}
	return ap, lv, true
}

// finalizeApproved marks a leave fully approved (approver = who) and upserts the
// year balance — the end of the approval chain.
func finalizeApproved(tx *gorm.DB, lv *models.Leave, approverID *string, now time.Time) error {
	lv.Status = "approved"
	lv.ApprovedAt = &now
	lv.ApprovedBy = approverID
	if err := tx.Save(lv).Error; err != nil {
		return err
	}
	return upsertLeaveYear(tx, lv)
}

// ApproveLeaveApproval processes one approver's pending task. If the chain is
// complete (this approver has role 'approval' and is at equal-or-higher
// authority than the requester) the leave is fully approved; otherwise the next
// 'approval'-role task is seeded up the hierarchy. No further approver found =>
// the request is finalized as approved.
func ApproveLeaveApproval(c *gin.Context) {
	ap, lv, ok := loadPendingStep(c)
	if !ok {
		return
	}
	if !authorizeApprover(c, &ap) {
		return
	}

	now := time.Now()
	curLevel := 0
	if ap.LLevel != nil {
		curLevel = *ap.LLevel
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		ap.Status = "approved"
		ap.ApprovedAt = &now
		if err := tx.Save(&ap).Error; err != nil {
			return err
		}

		reqLevel, hasReq := leaveflow.RequesterLevel(tx, lv.StaffID)
		complete := ap.RoleName == "approval" && hasReq && reqLevel >= curLevel
		if complete {
			return finalizeApproved(tx, &lv, ap.StaffID, now)
		}

		next := leaveflow.NextApproval(tx, lv.StaffID, curLevel)
		// No further approver, or it resolves back to the current actor/tier:
		// nothing left to escalate to, so finalize as approved.
		if next == nil || (ap.StaffID != nil && next.StaffID == *ap.StaffID && next.LevelType == curLevel) {
			return finalizeApproved(tx, &lv, ap.StaffID, now)
		}
		return insertApproval(tx, &lv, next)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to approve task"})
		return
	}

	enrichLeaves([]models.Leave{lv})
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"approval": ap, "leave": lv}})
}

// RejectLeaveApproval rejects one approver's task, which rejects the whole
// request (status + reason + who/when on the main leave).
func RejectLeaveApproval(c *gin.Context) {
	ap, lv, ok := loadPendingStep(c)
	if !ok {
		return
	}
	if !authorizeApprover(c, &ap) {
		return
	}

	var body struct {
		RejectReason string `json:"reject_reason"`
	}
	_ = c.ShouldBindJSON(&body)

	now := time.Now()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		ap.Status = "rejected"
		ap.ApprovedAt = &now
		if err := tx.Save(&ap).Error; err != nil {
			return err
		}
		lv.Status = "rejected"
		lv.RejectReason = body.RejectReason
		lv.ApprovedBy = ap.StaffID
		lv.ApprovedAt = &now
		return tx.Save(&lv).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to reject task"})
		return
	}

	enrichLeaves([]models.Leave{lv})
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"approval": ap, "leave": lv}})
}

// ---- leave_requester_files ----

var leaveFileSearch = []string{"comment", "staff_id", "added_by"}
var leaveFileFilter = []string{"leave_id", "task_id", "staff_id", "message_type", "status", "parent_id"}

func ListLeaveFiles(c *gin.Context) {
	listResource[models.LeaveRequesterFile](c, leaveFileSearch, leaveFileFilter, "id DESC")
}
func GetLeaveFile(c *gin.Context)    { getResource[models.LeaveRequesterFile](c, "id") }
func CreateLeaveFile(c *gin.Context) { createResource[models.LeaveRequesterFile](c) }
func UpdateLeaveFile(c *gin.Context) { updateResource[models.LeaveRequesterFile](c, "id") }
func DeleteLeaveFile(c *gin.Context) { deleteResource[models.LeaveRequesterFile](c, "id") }

// ---- leave_year ----

var leaveYearSearch = []string{"staff_id"}
var leaveYearFilter = []string{"staff_id", "leave_type_id", "l_year", "is_reset"}

func ListLeaveYears(c *gin.Context) {
	listResource[models.LeaveYear](c, leaveYearSearch, leaveYearFilter, "l_year DESC, staff_id ASC")
}
func GetLeaveYear(c *gin.Context) { getResource[models.LeaveYear](c, "id") }
