package shortid

import (
	"sync"
	"testing"
	"time"
)

func TestDuplicateKey(t *testing.T) {
	mutex := &sync.Mutex{}
	store := map[string]string{}
	total := 50000000

	wg := &sync.WaitGroup{}
	wg.Add(total)

	opt := Options{
		Number:        14,
		StartWithYear: true,
		EndWithHost:   false,
	}

	startTime := time.Now()

	for i := 0; i < total; i++ {
		key := Generate(opt)

		mutex.Lock()
		if _, ok := store[key]; ok {
			println("duplicate:", key)
		} else {
			store[key] = ""
		}
		mutex.Unlock()
		wg.Done()
	}
	wg.Wait()

	duration := int64(time.Since(startTime) / time.Millisecond)
	println("total count:", len(store))
	println("duration:", duration)
}

func TestSimple(t *testing.T) {
	opt := Options{
		Number:        14,
		StartWithYear: true,
		EndWithHost:   false,
	}

	for i := 0; i < 10; i++ {
		key := Generate(opt)
		println(key)
	}

}
