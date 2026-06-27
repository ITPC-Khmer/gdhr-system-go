package handlers

import (
	"net/http"
	"strings"

	"backend/database"
	"backend/models"
	"backend/services/hierarchy"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// This file exposes full CRUD (list + get + create + update + delete) for the
// four synced GDHR entities. List supports pagination, free-text search across
// a whitelist of columns, and exact-match filters on a whitelist of columns.
//
// The generic helpers below keep every entity's handlers to a single line each.

// paginate reads page/limit query params with sane bounds.
func paginate(c *gin.Context) (page, limit, offset int) {
	page = atoiDefault(c.DefaultQuery("page", "1"), 1)
	limit = atoiDefault(c.DefaultQuery("limit", "20"), 20)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 200 {
		limit = 20
	}
	return page, limit, (page - 1) * limit
}

func atoiDefault(s string, def int) int {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return def
		}
		n = n*10 + int(r-'0')
	}
	if s == "" {
		return def
	}
	return n
}

// listResource handles a paginated/searchable/filterable list for model T.
//   - searchCols: columns OR-matched with LIKE against ?search=
//   - filterCols: columns exact-matched against ?<col>=value (only when present)
//   - order:      default ORDER BY clause
// applyListFilters adds the free-text search and exact-match filters to q.
func applyListFilters(c *gin.Context, q *gorm.DB, searchCols, filterCols []string) *gorm.DB {
	if s := strings.TrimSpace(c.Query("search")); s != "" && len(searchCols) > 0 {
		like := "%" + s + "%"
		conds := make([]string, len(searchCols))
		args := make([]any, len(searchCols))
		for i, col := range searchCols {
			conds[i] = col + " LIKE ?"
			args[i] = like
		}
		q = q.Where(strings.Join(conds, " OR "), args...)
	}
	for _, col := range filterCols {
		if v := strings.TrimSpace(c.Query(col)); v != "" {
			q = q.Where(col+" = ?", v)
		}
	}
	return q
}

func listResource[T any](c *gin.Context, searchCols, filterCols []string, order string) {
	page, limit, offset := paginate(c)

	var model T
	q := applyListFilters(c, database.DB.Model(&model), searchCols, filterCols)

	var total int64
	q.Count(&total)

	rows := make([]T, 0, limit)
	if err := q.Order(order).Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rows, "total": total, "page": page, "limit": limit})
}

func getResource[T any](c *gin.Context, pkCol string) {
	var row T
	if err := database.DB.Where(pkCol+" = ?", c.Param("id")).First(&row).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": row})
}

func createResource[T any](c *gin.Context) {
	var row T
	if err := c.ShouldBindJSON(&row); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
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

// updateResource loads the existing row then binds the JSON body onto it, so
// only the fields present in the request are changed (partial update).
func updateResource[T any](c *gin.Context, pkCol string) {
	var row T
	if err := database.DB.Where(pkCol+" = ?", c.Param("id")).First(&row).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "record not found"})
		return
	}
	if err := c.ShouldBindJSON(&row); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := database.DB.Save(&row).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update record"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": row})
}

func deleteResource[T any](c *gin.Context, pkCol string) {
	if err := database.DB.Where(pkCol+" = ?", c.Param("id")).Delete(new(T)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete record"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "record deleted"})
}

// ---- institutes ----

var instituteSearch = []string{"id", "old_id", "name", "name_short", "source_table"}
var instituteFilter = []string{"parent_id", "level_type", "active", "source_table", "source_table_id"}

// instituteNamesFor returns a id->name map for the given (deduped, non-empty)
// institute ids, in a single query. Used to enrich rows that reference an
// institute by id (no foreign key).
func instituteNamesFor(ids []string) map[string]string {
	out := map[string]string{}
	if len(ids) == 0 {
		return out
	}
	var insts []models.Institute
	database.DB.Select("id", "name").Where("id IN ?", ids).Find(&insts)
	for _, in := range insts {
		out[in.ID] = in.Name
	}
	return out
}

func uniqueNonEmpty(values []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, v := range values {
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

// ListInstitutes is custom (not the generic helper) so it can resolve each
// row's parent_id to the parent institute's name (same table, no FK).
func ListInstitutes(c *gin.Context) {
	page, limit, offset := paginate(c)

	q := applyListFilters(c, database.DB.Model(&models.Institute{}), instituteSearch, instituteFilter)

	var total int64
	q.Count(&total)

	rows := make([]models.Institute, 0, limit)
	if err := q.Order("level_type ASC, name ASC").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch records"})
		return
	}

	ids := make([]string, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.ParentID)
	}
	names := instituteNamesFor(uniqueNonEmpty(ids))
	for i := range rows {
		rows[i].ParentName = names[rows[i].ParentID]
	}

	c.JSON(http.StatusOK, gin.H{"data": rows, "total": total, "page": page, "limit": limit})
}

func GetInstitute(c *gin.Context) { getResource[models.Institute](c, "id") }
func CreateInstitute(c *gin.Context) { createResource[models.Institute](c) }
func UpdateInstitute(c *gin.Context) { updateResource[models.Institute](c, "id") }
func DeleteInstitute(c *gin.Context) { deleteResource[models.Institute](c, "id") }

// ---- staffs ----

var staffSearch = []string{"uid", "staff_number", "surname_kh", "name_kh", "surname_en", "name_en", "phone", "email"}
var staffFilter = []string{
	"staff_type_id", "rank_id", "position_id", "other_position_id",
	"general_commissariat_id", "department_id", "office_id", "sector_id",
	"institute_id", "status_id", "gender",
}

// ListStaffs is custom so it can resolve each row's institute_id to the
// institute's name (no foreign key).
func ListStaffs(c *gin.Context) {
	page, limit, offset := paginate(c)

	q := applyListFilters(c, database.DB.Model(&models.Staff{}), staffSearch, staffFilter)

	var total int64
	q.Count(&total)

	rows := make([]models.Staff, 0, limit)
	if err := q.Order("name_kh").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
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

func GetStaff(c *gin.Context) { getResource[models.Staff](c, "uid") }

// CreateStaff/UpdateStaff are custom so they recompute the institute chain
// (institute_ids + institute_hierarchy) from institute_id on write.
func CreateStaff(c *gin.Context) {
	var staff models.Staff
	if err := c.ShouldBindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	staff.InstituteIDs, staff.InstituteHierarchy = hierarchy.Chain(database.DB, staff.InstituteID)
	if err := database.DB.Create(&staff).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"message": "a record with this key already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create record"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": staff})
}

func UpdateStaff(c *gin.Context) {
	var staff models.Staff
	if err := database.DB.Where("uid = ?", c.Param("id")).First(&staff).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "record not found"})
		return
	}
	if err := c.ShouldBindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	staff.InstituteIDs, staff.InstituteHierarchy = hierarchy.Chain(database.DB, staff.InstituteID)
	if err := database.DB.Save(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update record"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": staff})
}
func DeleteStaff(c *gin.Context) { deleteResource[models.Staff](c, "uid") }

// ---- ranks ----

var rankSearch = []string{"rank_name", "rank_name_short", "rank_name_en", "rank_name_short_en"}
var rankFilter = []string{"position_base_id", "rank_order", "active"}

func ListRanks(c *gin.Context)  { listResource[models.Rank](c, rankSearch, rankFilter, "rank_id") }
func GetRank(c *gin.Context)    { getResource[models.Rank](c, "rank_id") }
func CreateRank(c *gin.Context) { createResource[models.Rank](c) }
func UpdateRank(c *gin.Context) { updateResource[models.Rank](c, "rank_id") }
func DeleteRank(c *gin.Context) { deleteResource[models.Rank](c, "rank_id") }

// ---- positions ----

var positionSearch = []string{"position_name", "position_name_short", "description"}
var positionFilter = []string{"position_base_id", "rank_base_id"}

func ListPositions(c *gin.Context)  { listResource[models.Position](c, positionSearch, positionFilter, "position_id") }
func GetPosition(c *gin.Context)    { getResource[models.Position](c, "position_id") }
func CreatePosition(c *gin.Context) { createResource[models.Position](c) }
func UpdatePosition(c *gin.Context) { updateResource[models.Position](c, "position_id") }
func DeletePosition(c *gin.Context) { deleteResource[models.Position](c, "position_id") }

// ---- holidays ----

var holidaySearch = []string{"name"}
var holidayFilter = []string{"is_recurring"}

// authUserID returns the authenticated user's id from the gin context (set by
// the auth middleware), or nil when unavailable.
func authUserID(c *gin.Context) *uint {
	if v, ok := c.Get("userID"); ok {
		if id, ok := v.(uint); ok && id != 0 {
			return &id
		}
	}
	return nil
}

func ListHolidays(c *gin.Context) {
	listResource[models.Holiday](c, holidaySearch, holidayFilter, "date DESC")
}
func GetHoliday(c *gin.Context) { getResource[models.Holiday](c, "holiday_id") }

// CreateHoliday stamps created_by/updated_by from the authenticated user.
func CreateHoliday(c *gin.Context) {
	var h models.Holiday
	if err := c.ShouldBindJSON(&h); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	uid := authUserID(c)
	h.CreatedByID = uid
	h.UpdatedByID = uid
	if err := database.DB.Create(&h).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create record"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": h})
}

// UpdateHoliday loads the row, binds the (partial) body, and stamps updated_by.
func UpdateHoliday(c *gin.Context) {
	var h models.Holiday
	if err := database.DB.Where("holiday_id = ?", c.Param("id")).First(&h).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "record not found"})
		return
	}
	if err := c.ShouldBindJSON(&h); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	h.UpdatedByID = authUserID(c)
	if err := database.DB.Save(&h).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update record"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": h})
}

func DeleteHoliday(c *gin.Context) { deleteResource[models.Holiday](c, "holiday_id") }
