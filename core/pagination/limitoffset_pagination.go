package pagination

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LimitOffsetPagination[T any] struct {
	DefaultLimit int
	MaxLimit     int
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

func (p *LimitOffsetPagination[T]) Paginate(c *fiber.Ctx, db *gorm.DB) (*Response[T], error) {
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

	resp := &Response[T]{}
	var results []T

	var totalCount int64

	if err := db.Model(&results).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	countUint := uint(totalCount)
	resp.Count = &countUint

	if err := db.Limit(limit).Offset(offset).Find(&results).Error; err != nil {
		return nil, err
	}

	resp.Results = results
	resp.HasNext = (offset + limit) < int(totalCount)

	return resp, nil
}
