package mediamap

import (
	"timeline/internal/entity/dto/s3dto"
	"timeline/internal/infrastructure/models"
)

// Конвертация
func URLToDTO(urls []*models.ImageMeta) []*s3dto.FileURL {
	resp := make([]*s3dto.FileURL, 0, len(urls))
	for i := range urls {
		resp = append(resp, &s3dto.FileURL{URL: urls[i].URL, Type: urls[i].Type})
	}
	return resp
}
