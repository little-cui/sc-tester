package tests

import (
	"context"
	"fmt"
	"testing"
)

type Test1 struct {
	items map[int]*Item
	ch    chan *Item
	sch   chan []*Item
	ach   chan [NUM]*Item
}

type Test2 struct {
	items map[int]Item
	ch    chan Item
	sch   chan []Item
	ach   chan [NUM]Item
}

type Item struct {
	i int
}

func (i *Item) String() string {
	return fmt.Sprintf("{%d}", i.i)
}

const NUM = 1000

var t1 = &Test1{
	items: make(map[int]*Item, NUM),
	ch:    make(chan *Item, NUM),
	sch:   make(chan []*Item, NUM),
	ach:   make(chan [NUM]*Item, NUM),
}

var t2 = &Test2{
	items: make(map[int]Item, NUM),
	ch:    make(chan Item, NUM),
	sch:   make(chan []Item, NUM),
	ach:   make(chan [NUM]Item, NUM),
}

func BenchmarkWithDeferHandler1(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case a := <-t1.ch:
				a.i = a.i
			case b := <-t1.sch:
				b[0].i = b[0].i
				//fmt.Println("slice", len(b), cap(b), b)
			}
		}
	}()
	var (
		arr = make([]*Item, NUM)
		j   int
	)
	for i := 1; i < b.N; i++ {
		n := &Item{i: i}
		t1.items[i] = n
		t1.ch <- n
		if j < NUM {
			arr[j] = n
			j++
			continue
		}
		t1.sch <- arr[:j]
		arr = make([]*Item, NUM)
		j = 0
	}
	if j > 0 {
		t1.sch <- arr[:j]
	}
	b.ReportAllocs()
	cancel()
	// 2000000	       543 ns/op	      58 B/op	       1 allocs/op
}

func BenchmarkWithDeferHandler2(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case a := <-t2.ch:
				a.i = a.i
			case b := <-t2.sch:
				b[0].i = b[0].i
				fmt.Println("slice", len(b), cap(b), b)
			case c := <-t2.ach:
				c[0].i = c[0].i
				//fmt.Println("arr", len(c), cap(c), c)
			}
		}
	}()
	var (
		arr = [NUM]Item{}
		//arr = make([]Item, NUM)
		j int
	)
	for i := 1; i < b.N; i++ {
		n := Item{i: i}
		t2.items[i] = n
		t2.ch <- n
		if j < NUM {
			arr[j] = n
			j++
			continue
		}
		//t2.sch <- arr[:j] // BUG
		t2.ach <- arr
		//arr = make([]Item, NUM)
		arr = [NUM]Item{}
		j = 0
	}
	if j > 0 {
		//t2.sch <- arr[:j] // BUG
		t2.ach <- arr
	}
	b.ReportAllocs()
	cancel()
	// 2000000	       538 ns/op	      43 B/op	       0 allocs/op
}
