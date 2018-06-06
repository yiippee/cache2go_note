package main

import (
	"fmt"
	"time"

	"github.com/muesli/cache2go"
)

func main() {
	cache := cache2go.Cache("myCache", 100, 100)

	// The data loader gets called automatically whenever something
	// tries to retrieve a non-existing key from the cache.
	// 如果获取数据时，缓存未命中，则从这个函数中获取
	cache.SetDataLoader(func(key interface{}, args ...interface{}) *cache2go.CacheItem {
		// Apply some clever loading logic here, e.g. read values for
		// this key from database, network or file.
		val := "This is a test with key " + key.(string)

		// This helper method creates the cached item for us. Yay!
		item := cache2go.NewCacheItem(key, 0, val)
		return item
	})
	item, err := cache.Value("missGetCache")
	if err == nil {
		fmt.Println("load item:", item)
	}
	// This callback will be triggered every time a new item
	// gets added to the cache.
	cache.SetAddedItemCallback(func(entry *cache2go.CacheItem) {
		fmt.Println("Added:", entry.Key(), entry.Data(), entry.CreatedOn())
	})
	// This callback will be triggered every time an item
	// is about to be removed from the cache.
	cache.SetAboutToDeleteItemCallback(func(entry *cache2go.CacheItem) {
		fmt.Println("Deleting:", entry.Key(), entry.Data(), entry.CreatedOn())
	})

	// 测试 Add 带有修改功能吗？重新赋值就相当于修改了，会自动覆盖之前的value
	cache.Add("someKey", 0, "test update?")
	// Caching a new item will execute the AddedItem callback.
	cache.Add("someKey", 0, "This is a test!")

	// Let's retrieve the item from the cache
	res, err := cache.Value("someKey")
	if err == nil {
		fmt.Println("Found value in cache:", res.Data())
	} else {
		fmt.Println("Error retrieving value from cache:", err)
	}

	// Deleting the item will execute the AboutToDeleteItem callback.
	cache.Delete("someKey")

	// Caching a new item that expires in 3 seconds
	res = cache.Add("anotherKey", 3*time.Second, "This is another test")

	// This callback will be triggered when the item is about to expire
	res.SetAboutToExpireCallback(func(key interface{}) {
		fmt.Println("About to expire:", key.(string))
	})

	time.Sleep(5 * time.Second)
}
