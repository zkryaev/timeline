package datastore

type City struct {
	Name string `db:"name"`
	TZ   string `db:"tzid"`
}
