/*
 * Simple caching library with expiration capabilities
 *     Copyright (c) 2012, Radu Ioan Fericean
 *                   2013-2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package cache2go

import (
	"sync"
)

var (
	cache = make(map[string]*CacheTable)
	mutex sync.RWMutex
)

// Cache returns the existing cache table with given name or creates a new one
// if the table does not exist yet.
// 如果table存在，则返回，
// 如果不存在，则创建。
// cache table 是分组的，可以有多个分组，每个分组里面存具体的k-v
// 感觉类似于groupcache，所以可以分为 用户信息 俱乐部信息 等等所以可以分类的来管理
func Cache(table string, addChanSize, delChanSize int) *CacheTable {
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()
	var c1, c2 chan *CacheItem
	if !ok {
		mutex.Lock()
		t, ok = cache[table]
		// Double check whether the table exists or not.
		if !ok {
			if addChanSize > 0 {
				c1 = make(chan *CacheItem, addChanSize)
			}
			if delChanSize > 0 {
				c2 = make(chan *CacheItem, delChanSize)
			}
			t = &CacheTable{
				name:  table,
				items: make(map[interface{}]*CacheItem),
				addChan: c1,
				delChan: c2,
			}
			cache[table] = t
		}
		mutex.Unlock()
	}

	if t.addChan != nil {
		// 后台自动处理新增数据的流程
		go t.runAddFunc()
	}
	if t.delChan != nil {
		// 后台自动处理删数据的流程
		go t.runDelFunc()
	}

	return t
}
