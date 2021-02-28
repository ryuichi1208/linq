package main

import (
	"fmt"
	"strconv"
	"testing"
)

func TestExampleSuccess(t *testing.T) {
	l := readYaml("./test.yml")
	if len(l.Url) != 39 {
		t.Fatal("array reading error")
	}
}

func BenchmarkStrconvAppendFloat1(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.AppendFloat(nil, 3.1415926535, 'f', -1, 64)
	}
}

func BenchmarkAppend_AllocateEveryTime(b *testing.B) {
	base := []string{}
	b.ResetTimer()
	// Nはコマンド引数から与えられたベンチマーク時間から自動で計算される
	for i := 0; i < b.N; i++ {
		// 都度append
		base = append(base, fmt.Sprintf("no%d", i))
	}
}

func BenchmarkAppend_AllocateOnce(b *testing.B) {
	//最初に長さを決める
	base := make([]string, b.N)
	b.ResetTimer()
	// Nはコマンド引数から与えられたベンチマーク時間から自動で計算される
	for i := 0; i < b.N; i++ {
		base[i] = fmt.Sprintf("no%d", i)
	}
}
