package main_test

import (
    "strings"
    "unicode"
    "testing"
)

func SpaceMap(str string) string {
    return strings.Map(func(r rune) rune {
        if unicode.IsSpace(r) {
            return -1
        }
        return r
    }, str)
}

func SpaceFieldsJoin(str string) string {
    return strings.Join(strings.Fields(str), "")
}

func SpaceStringsBuilder(str string) string {
    var b strings.Builder
    b.Grow(len(str))
    for _, ch := range str {
        if !unicode.IsSpace(ch) {
            b.WriteRune(ch)
        }
    }
    return b.String()
}

func BenchmarkSpaceMap(b *testing.B) {
    for n := 0; n < b.N; n++ {
        SpaceMap(data)
    }
}

func BenchmarkSpaceFieldsJoin(b *testing.B) {
    for n := 0; n < b.N; n++ {
        SpaceFieldsJoin(data)
    }
}

func BenchmarkSpaceStringsBuilder(b *testing.B) {
    for n := 0; n < b.N; n++ {
        SpaceStringsBuilder(data)
    }
}
