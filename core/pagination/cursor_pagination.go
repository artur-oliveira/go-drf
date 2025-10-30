package pagination

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CursorPagination[T any] struct {
	DefaultLimit   int
	MaxLimit       int
	OrderByColumn  string
	OrderDirection string
}

func NewCursorPagination[T any](defaultLimit, maxLimit int, column, direction string) *CursorPagination[T] {
	if defaultLimit <= 0 {
		defaultLimit = 20
	}
	if maxLimit <= 0 {
		maxLimit = 100
	}
	if column == "" {
		column = "id"
	}
	if direction != "ASC" && direction != "DESC" {
		direction = "ASC"
	}
	return &CursorPagination[T]{
		DefaultLimit:   defaultLimit,
		MaxLimit:       maxLimit,
		OrderByColumn:  column,
		OrderDirection: strings.ToUpper(direction),
	}
}

func (p *CursorPagination[T]) Paginate(c *fiber.Ctx, db *gorm.DB) (*Response[T], error) {
	limit, _ := strconv.Atoi(c.Query("limit", strconv.Itoa(p.DefaultLimit)))
	cursor := c.Query("cursor", "")

	if limit > p.MaxLimit {
		limit = p.MaxLimit
	}
	if limit <= 0 {
		limit = p.DefaultLimit
	}

	direction := "asc"
	comparison := ">"
	if p.OrderDirection == "DESC" {
		direction = "desc"
		comparison = "<"
	}

	var results []T

	query := db.Order(p.OrderByColumn + " " + direction)

	if cursor != "" {
		query = query.Where(p.OrderByColumn+" "+comparison+" ?", cursor)
	}

	if err := query.Limit(limit + 1).Find(&results).Error; err != nil {
		return nil, err
	}

	resp := &Response[T]{Count: nil}

	if len(results) > limit {
		resp.HasNext = true
		resp.Results = results[:limit]
	} else {
		resp.HasNext = false
		resp.Results = results
	}

	return resp, nil
}
