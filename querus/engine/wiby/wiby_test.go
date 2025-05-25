package wiby_test

import (
	"fmt"
	"search/querus/engine"
	"search/querus/engine/wiby"
	"testing"
)

func TestSearch(t *testing.T) {
	moj := wiby.NewWiby(0.5)
	options := engine.RequestOptions{Query: "France", MaxResults: 24}
	data, err := moj.WebSearch(options)
	if err == nil {
		fmt.Println(data)
		fmt.Println(len(data))
	} else {
		fmt.Println(err)
	}
}
