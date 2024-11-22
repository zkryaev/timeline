package codemap

import (
	"timeline/internal/entity/dto"
	"timeline/internal/repository/models"
)

func ToModel(dtoCode *dto.VerifyCodeReq) *models.CodeInfo {
	return &models.CodeInfo{
		ID:    dtoCode.ID,
		Code:  dtoCode.Code,
		IsOrg: dtoCode.IsOrg,
	}
}
