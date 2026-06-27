package sync

import (
	"context"

	"backend/models"
	"backend/services/gdhr"
)

func mapInstitute(d gdhr.InstituteDTO) models.Institute {
	return models.Institute{
		ID:            d.ID,
		OldID:         d.OldID,
		ParentID:      d.ParentID,
		LevelType:     d.LevelType,
		OldParentID:   d.OldParentID,
		Name:          d.Name,
		NameShort:     d.NameShort,
		Active:        bool(d.Active),
		SourceTable:   d.SourceTable,
		SourceTableID: d.SourceTableID,
	}
}

func mapStaff(d gdhr.StaffDTO) models.Staff {
	return models.Staff{
		UID:                          d.UID,
		StaffTypeID:                  d.StaffTypeID,
		StaffNumber:                  d.StaffNumber,
		SurnameKh:                    d.SurnameKh,
		NameKh:                       d.NameKh,
		SurnameEn:                    d.SurnameEn,
		NameEn:                       d.NameEn,
		PlaceOfBirth:                 d.PlaceOfBirth,
		Nationality:                  d.Nationality,
		Address:                      d.Address,
		PhotoPath:                    d.PhotoPath,
		Gender:                       d.Gender,
		Phone:                        d.Phone,
		Email:                        d.Email,
		CityzenCardNumber:            d.CityzenCardNumber,
		RankID:                       d.RankID,
		RankNameShort:                d.RankNameShort,
		PositionID:                   d.PositionID,
		PositionNameShort:            d.PositionNameShort,
		OtherPositionID:              d.OtherPositionID,
		GeneralCommissariatID:        d.GeneralCommissariatID,
		GeneralCommissariatNameShort: d.GeneralCommissariatNameShort,
		DepartmentID:                 d.DepartmentID,
		DepartmentNameShort:          d.DepartmentNameShort,
		OfficeID:                     d.OfficeID,
		OfficeNameShort:              d.OfficeNameShort,
		SectorID:                     d.SectorID,
		SectorNameShort:              d.SectorNameShort,
		InstituteID:                  d.InstituteID,
		StatusID:                     d.StatusID,
		StatusName:                   d.StatusName,
	}
}

func mapRank(d gdhr.RankDTO) models.Rank {
	return models.Rank{
		RankID:          d.RankID,
		RankName:        d.RankName,
		RankNameShort:   d.RankNameShort,
		PositionBaseID:  d.PositionBaseID,
		RankOrder:       d.RankOrder,
		PromotePeriod:   d.PromotePeriod,
		RankNameEn:      d.RankNameEn,
		RankNameShortEn: d.RankNameShortEn,
		Active:          bool(d.Active),
	}
}

func mapPosition(d gdhr.PositionDTO) models.Position {
	return models.Position{
		PositionID:        d.PositionID,
		PositionName:      d.PositionName,
		Description:       d.Description,
		PositionNameShort: d.PositionNameShort,
		PositionBaseID:    d.PositionBaseID,
		RankBaseID:        d.RankBaseID,
	}
}

// CursorStatus reports the Redis cursor state for the paged entities on date.
func (s *Service) CursorStatus(ctx context.Context, date string) map[string]any {
	out := map[string]any{}
	for _, entity := range []string{"institutes", "staffs", "ranks", "positions"} {
		page, _ := NextPage(ctx, entity, date)
		done, _ := IsDone(ctx, entity, date)
		out[entity] = map[string]any{"next_page": page, "done": done}
	}
	return out
}
