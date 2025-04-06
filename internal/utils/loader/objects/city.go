package objects

type Timezone struct {
	TZ string `json:"tzid"`
}

type City struct {
	Name     string   `json:"name"`
	Timezone Timezone `json:"timezone"`
}

// like adapter pattern
type Cities struct {
	Arr      []City
	vidToInd *map[string]int
	lastVid  int
	cursor   int
}

func New(arr []City) Cities {
	dummymap := make(map[string]int, len(arr))
	fillmap(dummymap, arr)
	return Cities{
		Arr:      arr,
		vidToInd: &dummymap,
	}
}

func fillmap(dummymap map[string]int, arr []City) {
	for i := range arr {
		(dummymap)[arr[i].Name] = i
	}
}

func (c *Cities) AddCity(city City) {
	if _, ok := (*c.vidToInd)[city.Name]; ok {
		return
	}
	(*c.vidToInd)[city.Name] = c.lastVid
	c.Arr = append(c.Arr, City{Name: city.Name, Timezone: Timezone{TZ: city.Timezone.TZ}})
	c.lastVid++
}

func (c *Cities) GetCityTZ(name string) string {
	ind, ok := (*c.vidToInd)[name]
	if !ok {
		return ""
	}
	return c.Arr[ind].Timezone.TZ
}

func (c *Cities) Reset() {
	c.cursor = 0
}

func (c *Cities) Next() (City, bool) {
	if c.cursor >= len(c.Arr) {
		return City{}, false
	}
	defer func() { c.cursor++ }()
	return c.Arr[c.cursor], true
}
