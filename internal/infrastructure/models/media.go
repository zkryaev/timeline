package models

const (
	AVATAR   = 1
	SHOWCASE = 2
	BLOG     = 3
)

type ImageMeta struct {
	URL     string `db:"url"`
	Type    string `db:"type"`
	DomenID int
}
