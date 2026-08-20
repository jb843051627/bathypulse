package model

type Cursor struct {
	AfterID string `json:"after_id"`
	Limit   int    `json:"limit"`
}

func (c Cursor) Normalize() Cursor {
	if c.Limit <= 0 {
		c.Limit = 50
	}
	if c.Limit > 500 {
		c.Limit = 500
	}
	return c
}
