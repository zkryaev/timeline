package s3dto

import (
	"io"
	"timeline/internal/entity"
)

type FileURL struct {
	URL  string `json:"url"`
	Type string `json:"type"` // зарезервированное поле под блог в дальнейшем. Пока нету
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
	EntityID int    `json:"entity_id"`
	TData    entity.TokenData
}

type DeleteReq struct {
	Url    string
	Entity string
	TData  entity.TokenData
}
