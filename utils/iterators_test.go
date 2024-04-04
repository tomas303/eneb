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
			for {
				seq = seq + 1
				var x Result[int32, error]
				x.Value = seq
				channel <- x
			}
		})

	for i := 1; i < 1000000; i++ {
		if !iter.Next() {
			break
		}
		// fmt.Println(iter.current)
		if int32(i) != iter.current {
			t.Fatalf("expected %d but got %d", i, iter.current)
		}
	}
	close(iter.channel)
	elapsed := time.Since(start)
	fmt.Println("Test1 took ", elapsed)
}
