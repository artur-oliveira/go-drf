package pagination

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LimitOffsetPagination[T any] struct {
	DefaultLimit int
	MaxLimit     int

	limit  int
	offset int
}

func NewLimitOffsetPagination[T any](defaultLimit, maxLimit int) *LimitOffsetPagination[T] {
	if defaultLimit <= 0 {
		defaultLimit = 20
	}
	if maxLimit <= 0 {
		maxLimit = 100
	}
	return &LimitOffsetPagination[T]{
		DefaultLimit: defaultLimit,
		MaxLimit:     maxLimit,
	}
}

func (p *LimitOffsetPagination[T]) Bind(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", strconv.Itoa(p.DefaultLimit)))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > p.MaxLimit {
		limit = p.MaxLimit
	}
	if limit <= 0 {
		limit = p.DefaultLimit
	}
	if offset < 0 {
		offset = 0
	}

	p.limit = limit
	p.offset = offset
	return nil
}

func (p *LimitOffsetPagination[T]) Paginate(db *gorm.DB) (*Response[T], error) {
	resp := &Response[T]{}
	var results []T
	var totalCount int64

	if err := db.Model(&results).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	countUint := uint(totalCount)
	resp.Count = &countUint

	if err := db.Limit(p.limit).Offset(p.offset).Find(&results).Error; err != nil {
		return nil, err
	}

	resp.Results = results
	resp.HasNext = (p.offset + p.limit) < int(totalCount)

	return resp, nil
}
