package pagination

import (
	"fmt"
	"math"
	"strconv"
)

type Cursor struct {
	After  string `json:"after,omitempty"`
	Before string `json:"before,omitempty"`
	First  int    `json:"first,omitempty"`
	Last   int    `json:"last,omitempty"`
}

type PageInfo struct {
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
	StartCursor     string `json:"startCursor"`
	EndCursor       string `json:"endCursor"`
}

type PaginatedResult struct {
	Items    interface{} `json:"items"`
	PageInfo PageInfo    `json:"pageInfo"`
	Total    int         `json:"total"`
}

type OffsetPage struct {
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
	Offset  int `json:"-"`
}

func NewOffsetPage(page, perPage int) OffsetPage {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return OffsetPage{
		Page:    page,
		PerPage: perPage,
		Offset:  (page - 1) * perPage,
	}
}

func (p OffsetPage) TotalPages(total int) int {
	return int(math.Ceil(float64(total) / float64(p.PerPage)))
}

func EncodeCursor(id string, timestamp interface{}) string {
	return fmt.Sprintf("%s_%v", id, timestamp)
}

func DecodeCursor(cursor string) (id string, ok bool) {
	if len(cursor) == 0 {
		return "", false
	}
	return cursor, true
}

func ParseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	if v < 1 {
		return defaultVal
	}
	return v
}
