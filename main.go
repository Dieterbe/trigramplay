package main

import (
	"bufio"
	"errors"
	"fmt"
	index "github.com/raintank/raintank-metric/metric_tank/idx"
	"os"
	"runtime"
	//"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dgryski/trifles/repl"
)

func main() {

	var docs []string
	var ids []string
	var idx *index.Idx

	cBrute := func(args []string) error {
		if idx == nil {
			return errors.New("no index loaded")
		}

		if len(args) == 0 {
			return errors.New("missing argument")
		}

		patterns := args

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

	cIndex := func(args []string) error {
		runtime.GC()
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		fmt.Printf("ms.Alloc %d - ms.Sys %d\n", ms.Alloc, ms.Sys)
		if len(args) < 1 {
			return errors.New("missing argument")
		}
		fname := args[0]

		f, err := os.Open(fname)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanWords)

		if len(docs) != 0 {
			docs = docs[:0]
		}

		idx = index.New()

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
			id := idx.GetOrAdd(key)
			idx.AddRef(id)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("error during scan: ", err)
		}

		fmt.Printf("indexed %d documents in %s. index len %d \n", len(docs), time.Since(t0), idx.Len())
		runtime.GC()
		runtime.ReadMemStats(&ms)
		fmt.Printf("ms.Alloc %d - ms.Sys %d\n", ms.Alloc, ms.Sys)
		return nil
	}

	cPrint := func(args []string) error {
		for _, id := range ids {
			fmt.Println(id)
		}
		return nil
	}

	cPrune := func(args []string) error {
		if idx == nil {
			return errors.New("no index loaded")
		}

		if len(args) == 0 {
			return errors.New("missing argument")
		}

		pct, _ := strconv.Atoi(args[0])

		pruned := idx.Prune(float64(pct) / 100)
		fmt.Println("pruned", pruned, "at", pct)

		return nil
	}

	cSearch := func(args []string) error {
		if idx == nil {
			return errors.New("no index loaded")
		}

		if len(args) != 1 {
			return errors.New("need 1 query arg")
		}

		t0 := time.Now()
		// can't set to ids cause it's []Glob
		res := idx.Match(args[0])
		for _, r := range res {
			fmt.Println(r)
		}

		fmt.Println("found", len(res), "documents in", time.Since(t0))

		return nil
	}
	/*
				cTop := func(args []string) error {
					var freq []int
					for _, v := range idx {
						freq = append(freq, len(v))
					}

					sort.Ints(freq)

					for i := 0; i < 100; i++ {
						fmt.Println(freq[len(freq)-1-i])
					}
					return nil
				}

			cTrigram := func(args []string) error {
				if idx == nil {
					return errors.New("no index loaded")
				}

				var ts []trigram.T
				for _, f := range args {
					ts = trigram.Extract(f, ts)
				}

				for _, t := range ts {
					fmt.Printf("%q: %d\n", t, len(idx[t]))
				}

				return nil
			}

		cDelete := func(args []string) error {
			if idx == nil {
				return errors.New("no index loaded")
			}
			if len(args) < 1 {
				return errors.New("which id?")
			}

			id, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			t0 := time.Now()
			idx.Delete(args[0], trigram.DocID(id))
			fmt.Println("delete took", time.Since(t0))

			return nil
		}
	*/

	repl.Run("trigram> ",
		map[string]repl.Cmd{
			"brute": cBrute,
			//"delete":  cDelete,
			"index":  cIndex,
			"print":  cPrint,
			"prune":  cPrune,
			"search": cSearch,
			//	"top":     cTop,
			//	"trigram": cTrigram,
		},
	)
}
