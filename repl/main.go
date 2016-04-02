package main

import (
	"errors"
	"github.com/Dieterbe/trigramplay"
	"strconv"

	"github.com/dgryski/trifles/repl"
)

func main() {

	cBrute := func(args []string) error {
		if len(args) == 0 {
			return errors.New("missing pattern[s]")
		}
		return trigramplay.Brute(args)
	}

	cIndex := func(args []string) error {
		if len(args) < 1 {
			return errors.New("missing argument")
		}
		return trigramplay.Index(args[0])
	}

	cPrint := func(args []string) error {
		return trigramplay.PrintIds(args)
	}

	cPrune := func(args []string) error {
		if len(args) == 0 {
			return errors.New("missing argument")
		}
		pct, _ := strconv.Atoi(args[0])
		return trigramplay.Prune(pct)
	}

	cSearch := func(args []string) error {
		if len(args) != 1 {
			return errors.New("need 1 query arg")
		}
		return trigramplay.Search(args[0])

	}
	cTop := func(args []string) error {
		return trigramplay.Top()
	}

	/*

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
			"top":    cTop,
			//	"trigram": cTrigram,
		},
	)
}
