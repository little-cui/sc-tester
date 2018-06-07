/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package tests

import (
	"fmt"
	"testing"
)

func BenchmarkMapSlice1(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	m := make(map[string]int, b.N)
	for i := 0; i < b.N; i++ {
		m[fmt.Sprint(i)] = 0
	}
	b.StartTimer()
	for range m {
	}
	b.ReportAllocs()
	// 1000000	        20.6 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkMapSlice2(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	m := make([]int, 0, b.N)
	for i := 0; i < b.N; i++ {
		m = append(m, 0)
	}
	b.StartTimer()
	for range m {
	}
	b.ReportAllocs()
	// 1000000	         0.35 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkMapSlice3(b *testing.B) {
	b.StopTimer()
	b.N = 1000000
	m := make(map[string]int, b.N)
	for i := 0; i < b.N; i++ {
		m[fmt.Sprint(i)] = 0
	}
	b.StartTimer()
	s := make([]int, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	for range s {
	}
	b.ReportAllocs()
	// 1000000	        24.0 ns/op	       8 B/op	       0 allocs/op
}
