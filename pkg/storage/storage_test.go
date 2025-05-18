package storage

import (
	"fmt"
	"testing"
	"time"
)

func TestSyncPool(t *testing.T) {
	ch1 := make(chan *struct{}, 3)
	select {
	case x := <-ch1:
		t.Logf("x = %v", x)
	default:
		t.Logf("default")
	}
}

func TestTimer(t *testing.T) {
	timer := time.NewTimer(3 * time.Second)

	tm := <-timer.C
	fmt.Println("1", tm)

	timer.Reset(3 * time.Second)
	tm = <-timer.C
	fmt.Println("2", tm)

	timer.Stop()
}

func TestChan(t *testing.T) {
	ch := make(chan int, 2)

	fmt.Printf("1 ch %d - %d\n", len(ch), cap(ch))
	ch <- 1
	fmt.Printf("2 ch %d - %d\n", len(ch), cap(ch))
	ch <- 2

	fmt.Println("3")
	ch <- 3

	fmt.Println("4")
}
