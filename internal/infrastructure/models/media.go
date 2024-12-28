package models

const (
	AVATAR   = 1
	SHOWCASE = 2
	BLOG     = 3
)

type ImageMeta struct {
	DomenID int
	UUID    string `db:"uuid"`
	Type    int    `db:"type"`
}
