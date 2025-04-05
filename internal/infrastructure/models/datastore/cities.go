package datastore

import (
	"timeline/internal/utils/loader/objects"
)

type City struct {
	Name string `db:"name"`
	TZ   string `db:"tzid"`
}

// like adapter pattern
type Cities struct {
	Arr    []objects.City
	city   City
	cursor int
}

func (c *Cities) Reset() {
	c.city.Name = ""
	c.city.TZ = ""
	c.cursor = 0
}

func (c *Cities) Next() (City, bool) {
	if c.cursor >= len(c.Arr) {
		return City{}, false
	}
	c.city.Name = c.Arr[c.cursor].Name
	c.city.TZ = c.Arr[c.cursor].Timezone.TZ
	c.cursor++
	return c.city, true
}
