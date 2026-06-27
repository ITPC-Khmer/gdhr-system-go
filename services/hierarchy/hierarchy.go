// Package hierarchy resolves an institute's ancestor chain (self -> parent ->
// ... -> root) from the institutes table (parent_id links, no foreign key).
package hierarchy

import (
	"backend/models"

	"gorm.io/gorm"
)

const maxDepth = 50 // cycle / runaway guard

// Chain walks parent_id from instituteID up to the root, returning the id-chain
// and the matching name-chain (both starting with instituteID itself). Suited
// for one-off use (CRUD); it issues one small query per level.
func Chain(db *gorm.DB, instituteID string) ([]string, []string) {
	ids := []string{}
	names := []string{}
	seen := map[string]bool{}

	cur := instituteID
	for i := 0; cur != "" && !seen[cur] && i < maxDepth; i++ {
		seen[cur] = true
		var inst models.Institute
		if err := db.Select("id", "parent_id", "name").Where("id = ?", cur).First(&inst).Error; err != nil {
			break
		}
		ids = append(ids, inst.ID)
		names = append(names, inst.Name)
		cur = inst.ParentID
	}
	return ids, names
}

// Index is an in-memory snapshot of the institute tree for fast bulk chain
// resolution (used by the sync / backfill over many staff).
type Index struct {
	parent map[string]string
	name   map[string]string
}

// BuildIndex loads the whole institute tree (id, parent_id, name) into memory.
func BuildIndex(db *gorm.DB) (*Index, error) {
	var insts []models.Institute
	if err := db.Select("id", "parent_id", "name").Find(&insts).Error; err != nil {
		return nil, err
	}
	ix := &Index{
		parent: make(map[string]string, len(insts)),
		name:   make(map[string]string, len(insts)),
	}
	for _, in := range insts {
		ix.parent[in.ID] = in.ParentID
		ix.name[in.ID] = in.Name
	}
	return ix, nil
}

// Chain resolves the id-chain and name-chain for instituteID from the index.
func (ix *Index) Chain(instituteID string) ([]string, []string) {
	ids := []string{}
	names := []string{}
	seen := map[string]bool{}

	cur := instituteID
	for i := 0; cur != "" && !seen[cur] && i < maxDepth; i++ {
		nm, ok := ix.name[cur]
		if !ok {
			break
		}
		seen[cur] = true
		ids = append(ids, cur)
		names = append(names, nm)
		cur = ix.parent[cur]
	}
	return ids, names
}
