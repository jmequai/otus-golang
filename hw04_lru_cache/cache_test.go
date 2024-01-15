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

		c.Set("first", 1)
		c.Set("second", 2)
		c.Set("third", 3)

		value1, ok1 := c.Get("first")
		require.True(t, ok1)
		require.Equal(t, 1, value1)

		value2, ok2 := c.Get("second")
		require.True(t, ok2)
		require.Equal(t, 2, value2)

		value3, ok3 := c.Get("third")
		require.True(t, ok3)
		require.Equal(t, 3, value3)

		c.Set("fourth", 4)

		_, ok11 := c.Get("first")
		require.False(t, ok11)

		value22, ok22 := c.Get("second")
		require.True(t, ok22)
		require.Equal(t, 2, value22)

		value33, ok33 := c.Get("third")
		require.True(t, ok33)
		require.Equal(t, 3, value33)

		value44, ok44 := c.Get("fourth")
		require.True(t, ok44)
		require.Equal(t, 4, value44)
	})

	t.Run("move logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("first", 1)  // 1
		c.Set("second", 2) // 2,1
		c.Set("third", 3)  // 3,2,1

		value1, ok1 := c.Get("first") // 1,3,2
		require.True(t, ok1)
		require.Equal(t, 1, value1)

		value2, ok2 := c.Get("second") // 2,1,3
		require.True(t, ok2)
		require.Equal(t, 2, value2)

		value3, ok3 := c.Get("third") // 3,2,1
		require.True(t, ok3)
		require.Equal(t, 3, value3)

		c.Get("third")     // 3,2,1
		c.Get("second")    // 2,3,1
		c.Get("first")     // 1,2,3
		c.Set("first", 11) // 11,2,3
		c.Get("first")     // 11,2,3
		c.Set("fourth", 4) // 4,11,2

		value11, ok11 := c.Get("first")
		require.True(t, ok11)
		require.Equal(t, 11, value11)

		value22, ok22 := c.Get("second")
		require.True(t, ok22)
		require.Equal(t, 2, value22)

		_, ok33 := c.Get("third")
		require.False(t, ok33)

		value44, ok44 := c.Get("fourth")
		require.True(t, ok44)
		require.Equal(t, 4, value44)
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
