package main

import rw "github.com/aethiopicuschan/randomword-go"

func main() {
	req, err := rw.NewRequest(rw.WithNumber(10))
	if err != nil {
		panic(err)
	}
	words, err := req.Fetch()
	if err != nil {
		panic(err)
	}
	for _, word := range words {
		println(word)
	}
}
