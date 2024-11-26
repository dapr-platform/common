package core

import (
	"net/http"
)

// Entity 实体接口
type Entity interface {
	GetID() string
}

// Repository 数据仓储接口
type Repository[T Entity] interface {
	Query(r *http.Request) ([]T, error)
	QueryPage(r *http.Request) (*PageResult[T], error)
	GetByID(r *http.Request, id string) (T, error)
	Upsert(r *http.Request, entity T) error
	Delete(r *http.Request, id string) error
	BatchDelete(r *http.Request, ids []string) error
}

// PageResult 分页结果
type PageResult[T any] struct {
	Items    []T `json:"items"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
