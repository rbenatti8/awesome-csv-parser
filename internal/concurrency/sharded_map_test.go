package concurrency

import "testing"

func TestNewShardedMap(t *testing.T) {
	sm := NewShardedMap[int](10)
	if len(sm) != 10 {
		t.Errorf("Expected 10 shards, got %d", len(sm))
	}
}

func TestShardedMap_Set(t *testing.T) {
	sm := NewShardedMap[int](10)
	sm.Set("foo", 1)
	val, ok := sm.Get("foo")
	if !ok {
		t.Errorf("Expected to find key 'foo'")
	}
	if val != 1 {
		t.Errorf("Expected value of 1, got %d", val)
	}
}

func TestShardedMap_Delete(t *testing.T) {
	sm := NewShardedMap[int](10)
	sm.Set("foo", 1)
	sm.Delete("foo")
	_, ok := sm.Get("foo")
	if ok {
		t.Errorf("Expected to not find key 'foo'")
	}
}

func TestShardedMap_Get(t *testing.T) {
	sm := NewShardedMap[int](10)
	sm.Set("foo", 1)
	val, ok := sm.Get("foo")
	if !ok {
		t.Errorf("Expected to find key 'foo'")
	}
	if val != 1 {
		t.Errorf("Expected value of 1, got %d", val)
	}
}
