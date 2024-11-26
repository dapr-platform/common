package core

import (
	"context"
	"sync"
)

// HookType 定义hook类型
type HookType string

const (
	BeforeQuery       HookType = "before_query"
	AfterQuery        HookType = "after_query"
	BeforeUpsert      HookType = "before_upsert"
	AfterUpsert       HookType = "after_upsert"
	BeforeDelete      HookType = "before_delete"
	AfterDelete       HookType = "after_delete"
	BeforeBatchDelete HookType = "before_batch_delete"
	AfterBatchDelete  HookType = "after_batch_delete"
)

// HookFunc Hook函数定义
type HookFunc[T any] func(ctx context.Context, data T) (T, error)

// HookManager Hook管理器
type HookManager[T any] struct {
	mu    sync.RWMutex
	hooks map[HookType][]HookFunc[T]
}

// NewHookManager 创建Hook管理器
func NewHookManager[T any]() *HookManager[T] {
	return &HookManager[T]{
		hooks: make(map[HookType][]HookFunc[T]),
	}
}

// RegisterHook 注册Hook
func (hm *HookManager[T]) RegisterHook(hookType HookType, fn HookFunc[T]) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hm.hooks[hookType] = append(hm.hooks[hookType], fn)
}

// ExecuteHooks 执行Hooks
func (hm *HookManager[T]) ExecuteHooks(ctx context.Context, hookType HookType, data T) (T, error) {
	hm.mu.RLock()
	hooks := hm.hooks[hookType]
	hm.mu.RUnlock()

	var err error
	for _, hook := range hooks {
		data, err = hook(ctx, data)
		if err != nil {
			return data, err
		}
	}
	return data, nil
}
