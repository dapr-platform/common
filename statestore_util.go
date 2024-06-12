package common

import (
	"context"
	"encoding/json"
	"fmt"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/pkg/errors"
	"strconv"
	"sync"
	"time"
)

func CheckKeyInStateStore(ctx context.Context, client dapr.Client, stateStore, key string) (bool, error) {
	item, err := client.GetState(ctx, stateStore, key, nil)
	if err != nil {
		return false, err
	}

	if len(item.Value) == 0 { //没有这个key
		return false, nil
	}
	return true, nil
}

func DeleteKeyInStateStore(ctx context.Context, client dapr.Client, stateStore, key string) (err error) {
	return client.DeleteState(ctx, stateStore, key, make(map[string]string, 0))
}
func GetInStateStore(ctx context.Context, client dapr.Client, stateStore, key string) ([]byte, error) {
	item, err := client.GetState(context.Background(), stateStore, key, nil)
	if err != nil {
		return nil, err
	}
	if item == nil || len(item.Value) == 0 {
		return nil, nil
	}
	return item.Value, nil
}

func SaveInStateStore(ctx context.Context, client dapr.Client, stateStore, key string, data []byte, expires bool, ttl time.Duration) (err error) {
	ttlstr := "0"
	if expires {
		ttlstr = strconv.FormatInt(int64(ttl.Seconds()), 10)
	}
	item := &dapr.SetStateItem{
		Key: key,
		Etag: &dapr.ETag{
			Value: "1",
		},
		Metadata: map[string]string{
			"ttlInSeconds": ttlstr,
		},
		Value: data,
		Options: &dapr.StateOptions{
			Concurrency: dapr.StateConcurrencyLastWrite,
			Consistency: dapr.StateConsistencyStrong,
		},
	}

	return client.SaveBulkState(ctx, stateStore, item)

}

func BulkSaveInStateStore(ctx context.Context, client dapr.Client, stateStore string, key []string, data []string, expires bool, ttl time.Duration) (err error) {
	if len(key) != len(data) {
		return errors.New("key len !=data len")
	}
	ttlstr := "0"
	if expires {
		ttlstr = strconv.FormatFloat(ttl.Seconds(), 'f', 0, 64)
	}
	items := make([]*dapr.SetStateItem, 0)
	for i, k := range key {
		item := &dapr.SetStateItem{
			Key: k,
			Etag: &dapr.ETag{
				Value: "1",
			},
			Metadata: map[string]string{
				"ttlInSeconds": ttlstr,
			},
			Value: []byte(data[i]),
			Options: &dapr.StateOptions{
				Concurrency: dapr.StateConcurrencyLastWrite,
				Consistency: dapr.StateConsistencyStrong,
			},
		}
		items = append(items, item)
	}

	return client.SaveBulkState(ctx, stateStore, items...)

}

func argsToString(args ...any) string {
	if args == nil || len(args) == 0 {
		return ""
	}
	ret := ""
	for _, v := range args {
		ret += fmt.Sprintf("%v", v) + ","
	}
	if ret == "" {
		return ""
	}
	return ret[:len(ret)-1]
}
func DaprCacheGetGeneric[T any](ctx context.Context, client dapr.Client, dbFunc func(ctx context.Context, args ...any) (*T, error), key string, forceUpdate bool, expire bool, expireSeconds int, args ...any) (result *T, useStateStoreFlag bool, err error) {
	if forceUpdate {
		ret, err := dbFunc(ctx, args...)
		if err == nil {
			err = DaprCacheSet(ctx, client, key+argsToString(args...), ret, expire, expireSeconds)
		}
		return ret, false, err
	} else {
		//耗时操作，使用缓存
		ret, err := DaprCacheGet[T](client, key+argsToString(args...))
		if err != nil {
			Logger.Error("get Cache error " + err.Error())
			ret, err := dbFunc(ctx, args...)
			if err == nil {
				err = DaprCacheSet(ctx, client, key+argsToString(args...), ret, expire, expireSeconds)
			}
			return ret, false, err
		} else {
			if ret == nil {
				data, err := dbFunc(ctx, args...)
				if err != nil {
					return nil, false, errors.WithMessage(err, "db func error")
				}

				DaprCacheSet(ctx, client, key+argsToString(args...), data, expire, expireSeconds)
				return data, false, nil
			} else {
				return ret, true, nil
			}
		}
	}

}
func DaprCacheSet(ctx context.Context, client dapr.Client, key string, value interface{}, expire bool, expireSeconds int) error {

	buf, err := json.Marshal(value)
	if err != nil {
		Logger.Error("DparRedisCache Set marshal error")
		return err
	}
	return SaveInStateStore(ctx, client, DAPR_STATESTORE_NAME, key, buf, expire, time.Second*time.Duration(expireSeconds))

}

func DaprCacheGet[T any](client dapr.Client, key string) (*T, error) {

	item, err := client.GetState(context.Background(), DAPR_STATESTORE_NAME, key, nil)
	if err != nil {
		return nil, err
	}
	if item == nil || len(item.Value) == 0 {
		return nil, nil
	}
	var v T
	err = json.Unmarshal(item.Value, &v)
	if err != nil {
		Logger.Error("DaprCacheGet unmarshal error: " + string(item.Value))
		return nil, err
	}
	return &v, nil
}

type AutoRefreshCacher[T any] struct {
	Name          string
	Id            string
	ExpiredSecond int
	diffArgs      *sync.Map
	once          *sync.Once
	dbFunc        func(ctx context.Context, args ...any) (*T, error)
	key           string
	forceUpdate   bool
	expired       bool
}

func NewAutoRefreshCacher[T any](name string, expiredSecond int, expired bool, forceUpdate bool) *AutoRefreshCacher[T] {
	c := &AutoRefreshCacher[T]{
		Name:          name,
		Id:            NanoId(),
		ExpiredSecond: expiredSecond,
		diffArgs:      &sync.Map{},
		once:          &sync.Once{},
		forceUpdate:   forceUpdate,
	}
	return c
}
func (c *AutoRefreshCacher[T]) Invalid(client dapr.Client) {
	go func() {
		c.refresh(client)
	}()

	return
}
func (c *AutoRefreshCacher[T]) refresh(client dapr.Client) {
	c.diffArgs.Range(func(k, v interface{}) bool {
		_, _, err := DaprCacheGetGeneric[T](context.Background(), client, c.dbFunc, c.key, true, true, c.ExpiredSecond*10, v.([]any)...)
		if err != nil {
			Logger.Error("autoRefreshCacher error", err)
		}
		return true
	})
}

func (c *AutoRefreshCacher[T]) DaprCacheGetGeneric(ctx context.Context, client dapr.Client, dbFunc func(ctx context.Context, args ...any) (*T, error), key string, args ...any) (result *T, err error) {
	c.dbFunc = dbFunc
	c.key = key
	ret, _, err := DaprCacheGetGeneric[T](ctx, client, dbFunc, key, c.forceUpdate, c.expired, c.ExpiredSecond*10, args...)
	if err == nil {
		argKey := argsToString(args)
		c.diffArgs.Store(argKey, args)
	}

	c.once.Do(func() {

		go func(ctx context.Context, client dapr.Client, dbFunc func(ctx context.Context, args ...any) (*T, error), key string, forceUpdate bool, expire bool, expireSeconds int, args ...any) {
			if c.forceUpdate { //每次都更新，不需要定期刷新
				return
			}
			defaultDelay := time.Second * time.Duration(10)
			delay := defaultDelay
			for {

				time.Sleep(time.Second*time.Duration(c.ExpiredSecond) - delay)

				//当服务水平扩展多个时，需要保证刷新协程，只在一个服务实例中运行
				cacherKey := "AutoRefreshCacher:" + c.Name
				cacheExist, err := CheckKeyInStateStore(context.Background(), client, DAPR_STATESTORE_NAME, cacherKey)
				if err != nil {
					Logger.Error("check key in state error", err)
					continue
				}
				if !cacheExist { //第一次启动或超时
					SaveInStateStore(context.Background(), client, DAPR_STATESTORE_NAME, cacherKey, []byte(c.Id), true, time.Second*time.Duration(c.ExpiredSecond)+delay*2)
				}
				storeId, err := GetInStateStore(context.Background(), client, DAPR_STATESTORE_NAME, cacherKey)
				if err != nil {
					Logger.Error("getInStateStore err", err)
					continue
				}
				if string(storeId) != c.Id { //重新部署时，id不同，在超时前，会存在。
					Logger.Info(c.Id + " cacher " + c.Name + " process by other")
					continue
				}
				begin := time.Now()
				c.refresh(client)

				cost := time.Since(begin)
				Logger.Debug("autoCacher ", c.Name, " cost ", cost)
				if int(cost.Seconds()*2) > c.ExpiredSecond {
					Logger.Warning(c.Name+" one loop cost more time ", cost)
				}
				delay = time.Duration(cost.Seconds() * 2)
				if delay < defaultDelay {
					delay = defaultDelay
				}

			}
		}(ctx, client, dbFunc, key, false, true, c.ExpiredSecond, args...)

	})
	return ret, err
}
