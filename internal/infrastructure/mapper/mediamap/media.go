package mediamap

import (
	"timeline/internal/entity/dto/s3dto"
	"timeline/internal/infrastructure/models"
)

// Конвертация
func URLToDTO(URLs []*models.ImageMeta) []*s3dto.FileURL {
	resp := make([]*s3dto.FileURL, 0, len(URLs))
	for i := range URLs {
		resp = append(resp, &s3dto.FileURL{URL: URLs[i].URL, Type: URLs[i].Type})
	}
	return resp
}
