package handlers

import (
	"backend/models"

	"github.com/gin-gonic/gin"
)

// CRUD for leave_type and leave_roles. Both are plain reference tables, so they
// reuse the generic list/get/create/update/delete helpers from gdhr.go.

// ---- leave_type ----

var leaveTypeSearch = []string{"leave_key", "type_name", "type_name_s", "description"}
var leaveTypeFilter = []string{"is_reset"}

func ListLeaveTypes(c *gin.Context) {
	listResource[models.LeaveType](c, leaveTypeSearch, leaveTypeFilter, "id")
}
func GetLeaveType(c *gin.Context)    { getResource[models.LeaveType](c, "id") }
func CreateLeaveType(c *gin.Context) { createResource[models.LeaveType](c) }
func UpdateLeaveType(c *gin.Context) { updateResource[models.LeaveType](c, "id") }
func DeleteLeaveType(c *gin.Context) { deleteResource[models.LeaveType](c, "id") }

// ---- leave_roles ----

var leaveRoleSearch = []string{"leave_type", "name", "staff_type"}
var leaveRoleFilter = []string{"leave_type_id", "leave_type", "staff_type", "approve_level"}

func ListLeaveRoles(c *gin.Context) {
	listResource[models.LeaveRole](c, leaveRoleSearch, leaveRoleFilter, "leave_type_id ASC, min_duration ASC")
}
func GetLeaveRole(c *gin.Context)    { getResource[models.LeaveRole](c, "id") }
func CreateLeaveRole(c *gin.Context) { createResource[models.LeaveRole](c) }
func UpdateLeaveRole(c *gin.Context) { updateResource[models.LeaveRole](c, "id") }
func DeleteLeaveRole(c *gin.Context) { deleteResource[models.LeaveRole](c, "id") }