package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestIntIterator(t *testing.T) {
	start := time.Now()
	var seq int32 = 0
	var iter = NewIterator[int32](
		func(channel IteratorChannel[int32]) {
			seq = seq + 1
			var x Result[int32, error]
			x.Value = seq
			channel <- x
		})
	defer iter.Close()
	for i := 1; i < 1000000; i++ {
		if !iter.Next() {
			break
		}
		// fmt.Println(iter.current)
		if int32(i) != iter.current {
			t.Fatalf("expected %d but got %d", i, iter.current)
		}
	}
	elapsed := time.Since(start)
	t.Log(t.Name(), " took ", elapsed)
	iter = nil
}

func TestIntIterator2(t *testing.T) {
	start := time.Now()
	var seq int32 = 0
	var iter = NewIterator2[int32](
		func() (int32, bool) {
			seq = seq + 1
			if seq > 10 {
				return 0, false
			} else {
				return seq, true
			}
		})
	// for i := 1; i < 1000000; i++ {
	// 	if !iter.Next() {
	// 		break
	// 	}
	// 	// fmt.Println(iter.current)
	// 	if int32(i) != iter.current {
	// 		t.Fatalf("expected %d but got %d", i, iter.current)
	// 	}
	// }

	for iter.Next() {
		fmt.Println(iter.current)
	}
	error := iter.lastError
	if error != nil && error != IterEOF {
		t.Error(error)
	}

	elapsed := time.Since(start)
	t.Log(t.Name(), " took ", elapsed)
	iter = nil
}
