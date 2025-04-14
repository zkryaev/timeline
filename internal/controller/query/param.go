package query

type Param interface {
	load(val string) error
	isRequired() bool
	getName() string

	//GetType() string
	// GetInt() int
	// GetBool() bool
	// GetString() string
	// GetFloat32() float32
}
