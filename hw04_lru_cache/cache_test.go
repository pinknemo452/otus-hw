package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		requireFalse(t, ok)

		_, ok = c.Get("bbb")
		requireFalse(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		requireFalse(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		requireFalse(t, wasInCache)

		val, ok := c.Get("aaa")
		requireTrue(t, ok)
		requireEqual(t, 100, val)

		val, ok = c.Get("bbb")
		requireTrue(t, ok)
		requireEqual(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		requireTrue(t, wasInCache)

		val, ok = c.Get("aaa")
		requireTrue(t, ok)
		requireEqual(t, 300, val)

		val, ok = c.Get("ccc")
		requireFalse(t, ok)
		requireNil(t, val)
	})
	t.Run("evict oldest when full", func(t *testing.T) {
		c := NewCache(2)

		wasInCache := c.Set("a", 1)
		requireFalse(t, wasInCache)

		wasInCache = c.Set("b", 2)
		requireFalse(t, wasInCache)

		wasInCache = c.Set("c", 3)
		requireFalse(t, wasInCache)

		val, ok := c.Get("a")
		requireFalse(t, ok)
		requireNil(t, val)

		val, ok = c.Get("b")
		requireTrue(t, ok)
		requireEqual(t, 2, val)

		val, ok = c.Get("c")
		requireTrue(t, ok)
		requireEqual(t, 3, val)
	})

	t.Run("evict least recently used", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("a", 1)
		requireFalse(t, wasInCache)

		wasInCache = c.Set("b", 2)
		requireFalse(t, wasInCache)

		wasInCache = c.Set("c", 3)
		requireFalse(t, wasInCache)

		c.Get("a")
		c.Set("b", 22)
		c.Get("c")

		wasInCache = c.Set("d", 4)
		requireFalse(t, wasInCache)

		val, ok := c.Get("a")
		requireFalse(t, ok)
		requireNil(t, val)

		val, ok = c.Get("b")
		requireTrue(t, ok)
		requireEqual(t, 22, val)

		val, ok = c.Get("c")
		requireTrue(t, ok)
		requireEqual(t, 3, val)

		val, ok = c.Get("d")
		requireTrue(t, ok)
		requireEqual(t, 4, val)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
