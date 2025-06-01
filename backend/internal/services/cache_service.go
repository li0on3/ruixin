package services

import (
	"encoding/json"
	"sync"
	"time"
)

// CacheService 简单的内存缓存服务
type CacheService struct {
	mu    sync.RWMutex
	cache map[string]cacheItem
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewCacheService 创建缓存服务
func NewCacheService() *CacheService {
	cs := &CacheService{
		cache: make(map[string]cacheItem),
	}
	
	// 启动清理过期缓存的goroutine
	go cs.cleanupExpired()
	
	return cs
}

// Set 设置缓存
func (cs *CacheService) Set(key string, value interface{}, duration time.Duration) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	
	cs.cache[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(duration),
	}
	
	return nil
}

// Get 获取缓存
func (cs *CacheService) Get(key string) (interface{}, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	
	item, exists := cs.cache[key]
	if !exists {
		return nil, nil
	}
	
	// 检查是否过期
	if time.Now().After(item.expiration) {
		return nil, nil
	}
	
	return item.value, nil
}

// Delete 删除缓存
func (cs *CacheService) Delete(key string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	
	delete(cs.cache, key)
}

// Clear 清空所有缓存
func (cs *CacheService) Clear() {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	
	cs.cache = make(map[string]cacheItem)
}

// cleanupExpired 定期清理过期缓存
func (cs *CacheService) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		cs.mu.Lock()
		
		now := time.Now()
		for key, item := range cs.cache {
			if now.After(item.expiration) {
				delete(cs.cache, key)
			}
		}
		
		cs.mu.Unlock()
	}
}

// GetJSON 获取JSON格式的缓存
func (cs *CacheService) GetJSON(key string, target interface{}) error {
	data, err := cs.Get(key)
	if err != nil || data == nil {
		return err
	}
	
	// 如果数据已经是目标类型，直接返回
	if jsonData, ok := data.([]byte); ok {
		return json.Unmarshal(jsonData, target)
	}
	
	// 否则先序列化再反序列化
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(jsonData, target)
}

// SetJSON 设置JSON格式的缓存
func (cs *CacheService) SetJSON(key string, value interface{}, duration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	return cs.Set(key, jsonData, duration)
}