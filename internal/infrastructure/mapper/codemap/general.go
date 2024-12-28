package codemap

import (
	"timeline/internal/entity/dto/authdto"
	"timeline/internal/infrastructure/models"
)

func ToModel(dtoCode *authdto.VerifyCodeReq) *models.CodeInfo {
	return &models.CodeInfo{
		ID:    dtoCode.ID,
		Code:  dtoCode.Code,
		IsOrg: dtoCode.IsOrg,
	}
}
