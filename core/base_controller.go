package core

import (
	"net/http"
	"net/url"

	"github.com/dapr-platform/common"
)

// BaseController 基础控制器
type BaseController[T Entity] struct {
	Repository Repository[T]
	Hooks      *HookManager[[]T]
	SingleHooks *HookManager[T]
}

// NewBaseController 创建基础控制器
func NewBaseController[T Entity](repo Repository[T]) *BaseController[T] {
	return &BaseController[T]{
		Repository: repo,
		Hooks:      NewHookManager[[]T](),
		SingleHooks: NewHookManager[T](),
	}
}

func ParseQueryParams(r *http.Request) url.Values {
	return r.URL.Query()
}

// Query 查询处理
func (bc *BaseController[T]) Query(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 执行前置Hook
	var emptySlice []T
	result, err := bc.Hooks.ExecuteHooks(ctx, BeforeQuery, emptySlice)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	// 执行查询
	result, err = bc.Repository.Query(r)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	// 执行后置Hook
	result, err = bc.Hooks.ExecuteHooks(ctx, AfterQuery, result)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK.WithData(result))
}

// QueryPage 分页查询处理
func (bc *BaseController[T]) QueryPage(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("_page")
	pageSize := r.URL.Query().Get("_page_size")
	if page == "" || pageSize == "" {
		common.HttpResult(w, common.ErrParam.AppendMsg("page or pageSize is empty"))
		return
	}

	result, err := bc.Repository.QueryPage(r)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK.WithData(result))
}

// Upsert 创建或更新处理
func (bc *BaseController[T]) Upsert(w http.ResponseWriter, r *http.Request) {
	var entity T
	err := common.ReadRequestBody(r, &entity)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}

	ctx := r.Context()
	
	// 执行前置Hook
	entity, err = bc.SingleHooks.ExecuteHooks(ctx, BeforeUpsert, entity)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	err = bc.Repository.Upsert(r, entity)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	// 执行后置Hook
	entity, err = bc.SingleHooks.ExecuteHooks(ctx, AfterUpsert, entity)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpSuccess(w, common.OK.WithData(entity))
}

// Delete 删除处理
func (bc *BaseController[T]) Delete(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	
	// 获取实体用于Hook
	entity, err := bc.Repository.GetByID(r, id)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	// 执行前置Hook
	entity, err = bc.SingleHooks.ExecuteHooks(ctx, BeforeDelete, entity)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	err = bc.Repository.Delete(r, id)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	// 执行后置Hook
	entity, err = bc.SingleHooks.ExecuteHooks(ctx, AfterDelete, entity)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}

// BatchDelete 批量删除处理
func (bc *BaseController[T]) BatchDelete(w http.ResponseWriter, r *http.Request, ids []string) {
	if len(ids) == 0 {
		common.HttpResult(w, common.ErrParam.AppendMsg("len of ids is 0"))
		return
	}

	ctx := r.Context()
	var entities []T
	
	// 获取实体用于Hook
	for _, id := range ids {
		entity, err := bc.Repository.GetByID(r, id)
		if err != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
			return
		}
		entities = append(entities, entity)
	}

	// 执行前置Hook
	entities, err := bc.Hooks.ExecuteHooks(ctx, BeforeBatchDelete, entities)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	err = bc.Repository.BatchDelete(r, ids)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	// 执行后置Hook
	entities, err = bc.Hooks.ExecuteHooks(ctx, AfterBatchDelete, entities)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}
