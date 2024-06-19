package models

type Photo struct {
	ID  int    `db:"id"`
	URL string `db:"url"`
}
