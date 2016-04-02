package trigramplay

import (
	"bufio"
	"errors"
	"fmt"
	idx "github.com/raintank/raintank-metric/metric_tank/idx"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var docs []string
var ids []string
var loadedIndex *idx.Idx
var quiet bool

func Brute(patterns []string) error {
	if loadedIndex == nil {
		return errors.New("no index loaded")
	}

	if len(ids) != 0 {
		ids = ids[:0]
	}

	t0 := time.Now()
	for _, s := range docs {
		var mismatch = false
	search:
		for _, pat := range patterns {
			if !strings.Contains(s, pat) {
				mismatch = true
				break search
			}
		}

		if !mismatch {
			ids = append(ids, s)
		}
	}
	fmt.Println("found", len(ids), "documents in", time.Since(t0))
	return nil
}

func Index(fname string) error {
	var ms runtime.MemStats
	if !quiet {
		runtime.GC()
		runtime.ReadMemStats(&ms)
		fmt.Printf("ms.Alloc %d - ms.Sys %d\n", ms.Alloc, ms.Sys)
	}

	f, err := os.Open(fname)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	if len(docs) != 0 {
		docs = docs[:0]
	}

	loadedIndex = idx.New()

	t0 := time.Now()
	for scanner.Scan() {
		d := scanner.Text()
		orgId, err := strconv.Atoi(d)
		if err != nil {
			fmt.Println("bad orgid input", orgId)
		}
		scanner.Scan()
		key := scanner.Text()
		docs = append(docs, key)

		// add the trigrams
		id := loadedIndex.GetOrAdd(key)
		loadedIndex.AddRef(id)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("error during scan: ", err)
	}

	if !quiet {
		fmt.Printf("indexed %d documents in %s. index len %d \n", len(docs), time.Since(t0), loadedIndex.Len())
		runtime.GC()
		runtime.ReadMemStats(&ms)
		fmt.Printf("ms.Alloc %d - ms.Sys %d\n", ms.Alloc, ms.Sys)
	}
	return nil
}

func Get(key string) {
	id, found := loadedIndex.Get(key)
	fmt.Printf("id %d - found %t\n", id, found)
}

func GetOrAdd(key string) {
	id := loadedIndex.GetOrAdd(key)
	fmt.Printf("id %d\n", id)
}

func PrintIds(args []string) error {
	for _, id := range ids {
		fmt.Println(id)
	}
	return nil
}

func Prune(pct int) error {
	if loadedIndex == nil {
		return errors.New("no index loaded")
	}
	pruned := loadedIndex.Prune(float64(pct) / 100)
	if !quiet {
		fmt.Println("pruned", pruned, "at", pct)
	}

	return nil
}

func Search(query string) error {
	if loadedIndex == nil {
		return errors.New("no index loaded")
	}

	t0 := time.Now()
	// can't set to ids cause it's []Glob
	res := loadedIndex.Match(query)
	for _, r := range res {
		fmt.Println(r)
	}

	fmt.Println("found", len(res), "documents in", time.Since(t0))

	return nil

}

func Top(max int) {
	var freq []int
	for _, v := range loadedIndex.Pathidx {
		freq = append(freq, len(v))
	}

	sort.Ints(freq)
	if len(freq) < max {
		max = len(freq)
	}

	for i := 0; i < max; i++ {
		fmt.Println(freq[len(freq)-1-i])
	}
}

func Show() {
	for t, postList := range loadedIndex.Pathidx {
		if postList == nil {
			fmt.Printf("T %s - nil -> all docs. ids:")
			postList = loadedIndex.Pathidx[0xFFFFFFFF]
		} else {
			fmt.Printf("T %s - %d docs. ids:", t, len(postList))
		}
		for _, i := range postList {
			fmt.Printf(" %d", i)
		}
		fmt.Println()
		for _, id := range postList {
			fmt.Println(" ", loadedIndex.GetById(idx.MetricID(id)))
		}
	}
}
