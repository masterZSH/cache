package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	foo = 1
)

func TestAdd(t *testing.T) {
	k, v := 1, 2
	cache := DefaultCache()
	duration := 5 * time.Second
	item, result := cache.Add(
		k,
		v,
		duration,
		func() {
			fmt.Println("add success")
		},
	)
	assert.NotEmpty(t, item)
	assert.Equal(t, true, result)

	item, result = cache.Add(
		k,
		v,
		duration,
		func() {
			fmt.Println("add success")
		},
	)
	assert.NotEmpty(t, item)
	assert.Equal(t, false, result)
	assert.Equal(t, v, item)

	// override test
	oc := New(
		&CacheConfig{
			AutoOverride: true,
		},
	)
	item, result = oc.Add(
		k,
		v,
		duration,
		func() {
			fmt.Println("add success")
		},
	)
	assert.NotEmpty(t, item)
	assert.Equal(t, true, result)

	item, result = oc.Add(
		k,
		v,
		duration,
		func() {
			fmt.Println("add success")
		},
	)
	assert.NotEmpty(t, item)
	assert.Equal(t, true, result)

}

func TestGet(t *testing.T) {
	k, v := 1, 2
	c := DefaultCache()
	c.Add(k, v, 5*time.Second, func() {
		foo = 2
	})
	time.Sleep(6 * time.Second)

	// test callback
	assert.Equal(t, 2, foo)

	c.Add(k, v, 5*time.Second, func() {
		foo = 2
	})

	actual, ok := c.Get(k)
	assert.Equal(t, ok, true)
	assert.Equal(t, v, actual)
}
