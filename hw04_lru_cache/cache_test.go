package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 300, val)

		_, ok = c.Get("ddd")
		require.False(t, ok)

		c.Clear()

		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Equal(t, nil, val)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Equal(t, nil, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Equal(t, nil, val)

		wasInCache := c.Set("bbb", 500)
		require.False(t, wasInCache)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 500, val)
	})
}

func TestCacheEjection(t *testing.T) {
	t.Run("simple ejection logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		val, ok = c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 300, val)

		_, ok = c.Get("ddd")
		require.False(t, ok)

		wasInCache := c.Set("ddd", 400)
		require.False(t, wasInCache)

		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Equal(t, nil, val)

		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 400, val)
	})

	t.Run("ejection logic with old values", func(t *testing.T) {
		c := NewCache(10)

		for i := 0; i < 10; i++ {
			wasInCache := c.Set(Key("N"+strconv.Itoa(i)), 100*i)
			require.False(t, wasInCache)
		}

		val, ok := c.Get("N0")
		require.True(t, ok)
		require.Equal(t, 0, val)

		val, ok = c.Get("N9")
		require.True(t, ok)
		require.Equal(t, 900, val)

		val, ok = c.Get("N2")
		require.True(t, ok)
		require.Equal(t, 200, val)

		for i := 10; i < 13; i++ {
			wasInCache := c.Set(Key("N"+strconv.Itoa(i)), 100*i)
			require.False(t, wasInCache)
		}

		val, ok = c.Get("N0")
		require.True(t, ok)
		require.Equal(t, 0, val)

		val, ok = c.Get("N1")
		require.False(t, ok)
		require.Equal(t, nil, val)

		val, ok = c.Get("N2")
		require.True(t, ok)
		require.Equal(t, 200, val)

		val, ok = c.Get("N3")
		require.False(t, ok)
		require.Equal(t, nil, val)

		val, ok = c.Get("N4")
		require.False(t, ok)
		require.Equal(t, nil, val)

		val, ok = c.Get("N10")
		require.True(t, ok)
		require.Equal(t, 1000, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	// t.Skip() // Remove me if task with asterisk completed.

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
