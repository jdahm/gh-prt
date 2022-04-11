package main

import (
	"fmt"
	"testing"
)

type QueryMocker struct {
	page int
}

func (q *QueryMocker) Next() map[string]int {
	query := make(map[string]int)
	if q.page == 0 {
		query["bo3q4t8"] = 0
	} else {
		query["abcdefg"] = 1
		query["tuvwxyz"] = 4
	}
	q.page++
	return query
}

func (q *QueryMocker) AtEnd() bool {
	return q.page > 1
}

func TestFindMatchesNoMatch(t *testing.T) {
	candidates := []string{"ab cdefg", "tuvwxyz", "bo3q4t8"}
	results := findMatches("BWX", candidates)
	if len(results) != 0 {
		t.Error("Found non-result")
	}
}

func TestFindMatchesSingle(t *testing.T) {
	candidates := []string{"ab cdefg", "tuvwxyz", "bo3q4t8"}
	results := findMatches("ABC", candidates)
	if len(results) != 1 {
		t.Error("Did not find single result")
	} else if results[0] != "ab cdefg" {
		t.Error("Did not find correct result")
	}
}

func TestFindMatchingPRs(t *testing.T) {
	q := &QueryMocker{}
	results := FindMatchingPRs(q, "abc")
	fmt.Println(results)
	if len(results) != 1 {
		t.Error("Did not find result")
	}
	v, ok := results["abcdefg"]
	if !ok {
		t.Error("Correct result not in map")
	}
	if v != 1 {
		t.Error("Value is incorrect")
	}
}
