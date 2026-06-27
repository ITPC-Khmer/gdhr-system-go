package handlers

import (
	"net/http"
	"strings"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
)

// CRUD for staff_institute_roles — the staff<->institute management mapping
// (admin/moderator/approval) that the leave-approval routing reads. List
// resolves each row's institute_id to the institute name (no FK).

var staffInstituteRoleSearch = []string{"staff_id", "institute_id", "role"}
var staffInstituteRoleFilter = []string{"staff_id", "institute_id", "role"}

func ListStaffInstituteRoles(c *gin.Context) {
	page, limit, offset := paginate(c)

	q := applyListFilters(c, database.DB.Model(&models.StaffInstituteRole{}), staffInstituteRoleSearch, staffInstituteRoleFilter)

	var total int64
	q.Count(&total)

	rows := make([]models.StaffInstituteRole, 0, limit)
	if err := q.Order("institute_id ASC, role ASC").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch records"})
		return
	}

	ids := make([]string, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.InstituteID)
	}
	names := instituteNamesFor(uniqueNonEmpty(ids))
	for i := range rows {
		rows[i].InstituteName = names[rows[i].InstituteID]
	}

	c.JSON(http.StatusOK, gin.H{"data": rows, "total": total, "page": page, "limit": limit})
}

func GetStaffInstituteRole(c *gin.Context) { getResource[models.StaffInstituteRole](c, "id") }

// CreateStaffInstituteRole enforces the business rule that an institute may have
// at most ONE moderator and at most ONE approval (admin may be many).
func CreateStaffInstituteRole(c *gin.Context) {
	var row models.StaffInstituteRole
	if err := c.ShouldBindJSON(&row); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if row.Role == "moderator" || row.Role == "approval" {
		var count int64
		database.DB.Model(&models.StaffInstituteRole{}).
			Where("institute_id = ? AND role = ?", row.InstituteID, row.Role).
			Count(&count)
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"message": "institute already has a " + row.Role})
			return
		}
	}

	if err := database.DB.Create(&row).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"message": "a record with this key already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create record"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": row})
}

// UpdateStaffInstituteRole applies the same one-moderator/one-approval-per-
// institute rule as create, ignoring the row being updated.
func UpdateStaffInstituteRole(c *gin.Context) {
	var row models.StaffInstituteRole
	if err := database.DB.Where("id = ?", c.Param("id")).First(&row).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "record not found"})
		return
	}
	if err := c.ShouldBindJSON(&row); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if row.Role == "moderator" || row.Role == "approval" {
		var count int64
		database.DB.Model(&models.StaffInstituteRole{}).
			Where("institute_id = ? AND role = ? AND id <> ?", row.InstituteID, row.Role, row.ID).
			Count(&count)
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"message": "institute already has a " + row.Role})
			return
		}
	}

	if err := database.DB.Save(&row).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"message": "a record with this key already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update record"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": row})
}

func DeleteStaffInstituteRole(c *gin.Context) { deleteResource[models.StaffInstituteRole](c, "id") }
