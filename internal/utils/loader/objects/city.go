package objects

type Timezone struct {
	TZ string `json:"tzid"`
}

type City struct {
	Name     string   `json:"name"`
	Timezone Timezone `json:"timezone"`
}
