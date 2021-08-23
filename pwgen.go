// Correct Horse Battery Staple Password Generator
//
// https://xkcd.com/936/
//
// For RSI reasons, this is designed to avoid any characters that require a
// shift.
package main

import (
	"crypto/rand"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

// EFF's New Wordlists for Random Passpharses.
//
// https://www.eff.org/deeplinks/2016/07/new-wordlists-random-passphrases
var (
	effWordLists = map[string]string{
		"large": "https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt",
		"short": "https://www.eff.org/files/2016/09/08/eff_short_wordlist_1.txt",
	}
	seperatorList = []string{
		".",
		"-",
		" ",
		"=",
		"/",
		",",
	}
)

type generator func(entropy *big.Int) string

func joinGenerator(partGens []generator, sepGen generator) generator {
	return func(entropy *big.Int) string {
		sep := sepGen(entropy)
		var parts []string
		for _, partGen := range partGens {
			parts = append(parts, partGen(entropy))
		}
		return strings.Join(parts, sep)
	}
}

func makeStringListGenerator(items []string) generator {
	// doesn't actually pluck. maybe it should remove the choice from the word
	// list.
	length := big.NewInt(int64(len(items)))
	itemSet := NewStringSet(items...)
	items = itemSet.Sorted()
	return func(entropy *big.Int) string {
		n, err := rand.Int(rand.Reader, length)
		if err != nil {
			log.Fatal(err)
		}
		entropy.Mul(entropy, length)
		return items[n.Int64()]
	}
}

func makeWordGenerator() generator {
	resp, err := http.Get(effWordLists["short"])
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	r := csv.NewReader(resp.Body)
	r.Comma = '\t'

	var words []string
	for rec, err := r.Read(); err != io.EOF; rec, err = r.Read() {
		if err != nil {
			log.Fatal(err)
		}
		words = append(words, rec[1])
	}
	return makeStringListGenerator(words)
}

func makeNumGenerator() generator {
	var nums []string
	for i := 0; i < 100; i++ {
		nums = append(nums, strconv.Itoa(i))
	}
	return makeStringListGenerator(nums)
}

func wrapGenerator(bodyGen, quoteGen generator) generator {
	return func(entropy *big.Int) string {
		body := bodyGen(entropy)
		quote := quoteGen(entropy)
		return quote + body + quote
	}
}

var (
	numIters = flag.Int("num_iters", 10, "")
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	sepGen := makeStringListGenerator(seperatorList)
	wordGen := makeWordGenerator()
	numGen := makeNumGenerator()

	pwGen := wrapGenerator(
		joinGenerator([]generator{
			numGen,
			wordGen,
			wordGen,
			wordGen,
			numGen,
		}, sepGen),
		sepGen,
	)

	log.Printf("generating %d passwords", *numIters)
	for i := 0; i < *numIters; i++ {
		entropy := big.NewInt(1)
		pw := pwGen(entropy)
		if entropy.Int64() < 40 {
			log.Fatalf("entropy too low: %d", entropy)
		}
		fmt.Printf("password=%q, entropy=%d\n", pw, entropy.BitLen())
	}

}

type StringSet map[string]struct{}

func NewStringSet(items ...string) StringSet {
	out := make(StringSet)
	for _, item := range items {
		out[item] = struct{}{}
	}
	return out
}

func (ss StringSet) Sorted() []string {
	items := make([]string, 0)
	for item, _ := range ss {
		items = append(items, item)
	}
	return items
}
