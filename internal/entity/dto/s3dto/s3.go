package s3dto

import "io"

type FileURL struct {
	URL  string `json:"url"`
	Type int    `json:"type,omitempty"` // зарезервированное поле под блог в дальнейшем. Пока нету
}

type CreateFileDTO struct {
	DomenInfo
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	Reader io.Reader
}

type File struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Bytes []byte `json:"file"`
}

type DomenInfo struct {
	Entity   string `json:"entity"`
	EntityID int    `json:"entityID"`
	Type     string `json:"type"`
}
