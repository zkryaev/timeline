package models

type SearchParams struct {
	Page   int
	Limit  int
	Offset int
	Name   string
	Type   string
}

type AreaParams struct {
	Left  Coordinate
	Right Coordinate
}
