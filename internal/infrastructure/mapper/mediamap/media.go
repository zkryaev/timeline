package mediamap

import (
	"timeline/internal/entity/dto/s3dto"
)

// Конвертация
func ImageUUIDToURL(ImagesURL ...string) []*s3dto.FileURL {
	resp := make([]*s3dto.FileURL, 0, len(ImagesURL))
	for i := range ImagesURL {
		resp = append(resp, &s3dto.FileURL{URL: ImagesURL[i]})
	}
	return resp
}
