package validation

var (
	ImageTypes = map[string]struct{}{
		"image/jpeg": {},
		"image/png":  {},
		"image/webp": {},
		"image/gif":  {},
		"image/bmp":  {},
		"image/tiff": {},
	}
)

func IsImage(ContentType string) bool {
	if _, ok := ImageTypes[ContentType]; !ok {
		return false
	}
	return true
}
