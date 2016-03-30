package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dgryski/go-trigram"

	"github.com/dgryski/trifles/repl"
)

func main() {

	var docs []string
	var ids []trigram.DocID
	var idx trigram.Index

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
		for i, s := range docs {
			var mismatch = false
		search:
			for _, pat := range patterns {
				if !strings.Contains(s, pat) {
					mismatch = true
					break search
				}
			}

			if !mismatch {
				ids = append(ids, trigram.DocID(i))
			}
		}
		fmt.Println("found", len(ids), "documents in", time.Since(t0))

		return nil
	}

	cFilter := func(args []string) error {
		if idx == nil {
			return errors.New("no index loaded")
		}

		var ts []trigram.T
		for _, f := range args {
			ts = trigram.Extract(f, ts)
		}

		t0 := time.Now()
		ids = idx.Filter(ids, ts)

		fmt.Println("filtered", len(ids), "documents in", time.Since(t0))
		return nil
	}

	cIndex := func(args []string) error {
		if len(args) < 1 {
			return errors.New("missing argument")
		}
		fname := args[0]

		f, err := os.Open(fname)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(f)

		if len(docs) != 0 {
			docs = docs[:0]
		}

		idx = trigram.NewIndex(nil)

		t0 := time.Now()
		for scanner.Scan() {
			d := scanner.Text()
			docs = append(docs, d)

			// add the trigrams
			idx.Add(d)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("error during scan: ", err)
		}

		fmt.Printf("indexed %d documents in %s\n", len(docs), time.Since(t0))
		return nil
	}

	cPrint := func(args []string) error {
		for _, id := range ids {
			fmt.Printf("%05d: %q\n", id, docs[id])
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

		var ts []trigram.T
		for _, f := range args {
			ts = trigram.Extract(f, ts)
		}

		t0 := time.Now()
		ids = idx.QueryTrigrams(ts)

		fmt.Println("found", len(ids), "documents in", time.Since(t0))

		return nil
	}

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

	repl.Run("trigram> ",
		map[string]repl.Cmd{
			"brute":   cBrute,
			"delete":  cDelete,
			"filter":  cFilter,
			"index":   cIndex,
			"print":   cPrint,
			"prune":   cPrune,
			"search":  cSearch,
			"top":     cTop,
			"trigram": cTrigram,
		},
	)
}
