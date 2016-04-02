package trigramplay

import (
	"testing"
)

func init() {
	quiet = true
}

func search(q string) {
	res := loadedIndex.Match(q)
	if len(res) < 1 {
		panic("no matches??")
	}
}

func BenchmarkLiteral(b *testing.B) {
	Index("prod.txt")
	Prune(20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search("litmus.dieter_plaetinck_be.paris.ping.max")
	}
}
func BenchmarkStarFirst(b *testing.B) {
	Index("prod.txt")
	Prune(20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search("*.dieter_plaetinck_be.paris.ping.max")
	}
}
func BenchmarkStarLast(b *testing.B) {
	Index("prod.txt")
	Prune(20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search("litmus.dieter_plaetinck_be.paris.ping.*")
	}
}
func BenchmarkStarSecondLast(b *testing.B) {
	Index("prod.txt")
	Prune(20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search("litmus.dieter_plaetinck_be.paris.*.ok_state")
	}
}
func BenchmarkStarThirdLast(b *testing.B) {
	Index("prod.txt")
	Prune(20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search("litmus.dieter_plaetinck_be.*.ping.ok_state")
	}
}
func BenchmarkStarLastAndSecondLast(b *testing.B) {
	Index("prod.txt")
	Prune(20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search("litmus.dieter_plaetinck_be.paris.*.*")
	}
}
func BenchmarkStarSecondLastAndthirdLast(b *testing.B) {
	Index("prod.txt")
	Prune(20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search("litmus.dieter_plaetinck_be.*.*.ok_state")
	}
}
func BenchmarkStarLastAndSecondLastAndThirdLast(b *testing.B) {
	Index("prod.txt")
	Prune(20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search("litmus.dieter_plaetinck_be.*.*.*")
	}
}

func BenchmarkStarSecondAndLastAndSecondLastAndThirdLast(b *testing.B) {
	Index("prod.txt")
	Prune(20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search("litmus.*.*.*.*")
	}
}
func BenchmarkStarEverywhere(b *testing.B) {
	Index("prod.txt")
	Prune(20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search("*.*.*.*.*")
	}
}
