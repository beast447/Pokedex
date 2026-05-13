
package internal

import (
	"time"
	"testing"
	"fmt"
)


func CacheTest(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
		key: "www.example.com/api",
		val: []byte("lotsadata"),
		},
	
	{
		key: "http://anotherexample.com",
		val: []byte("moredatafrombytes"),
	},
	{
		key: "https://Thisisawebsite.com/api/go",
		val: []byte("Thisisthedatafromthewebsite"),
	},
}

for i, c := range cases {
	t.Run(fmt.Sprintf("Test case %v", i), func (t *testing.T) {
		cache := NewCache(interval)
		cache.Add(c.key, c.val)
		val, ok := cache.Get(c.key)
			if !ok{
				t.Errorf("expected to find a key")
				return
			}
			if string(val) != string(c.val){
				t.Errorf("expected to find value")
			}
		})
}

}
