package query

type Param interface {
	load(val string) error
	isRequired() bool
	getName() string
}
