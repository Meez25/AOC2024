package main

import "testing"

func BenchmarkArrayIteration(b *testing.B) {
	arr := make([]int, 1000000)
	for i := 0; i < b.N; i++ {
		for range arr {
		}
	}
}

func BenchmarkMapIteration(b *testing.B) {
	m := make(map[int]int, 1000000)
	for i := 0; i < b.N; i++ {
		for _ = range m {
		}
	}
}
