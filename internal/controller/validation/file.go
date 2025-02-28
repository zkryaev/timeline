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

func IsImage(contentType string) bool {
	if _, ok := ImageTypes[contentType]; !ok {
		return false
	}
	return true
}
