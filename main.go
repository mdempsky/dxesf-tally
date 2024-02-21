// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The dxesf-tally command computes the winners of a DxE SF Bay Area
// core election. It accepts ballots on standard input in CSV format
// and prints a report to standard output.
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
)

func main() {
	rows, err := csv.NewReader(os.Stdin).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	header, rows := rows[0], rows[1:]
	_ = header

	if false {
		names := header[2:]
		fmt.Println(names)
	}

	ballots := make([]ballot, len(rows))
	for i, row := range rows {
		ballot := &ballots[i]
		ballot.invWeight.SetInt64(1)
		for j := range ballot.scores {
			score, err := strconv.Atoi(row[2+j])
			if err != nil {
				log.Fatal(err)
			}
			ballot.scores[j].SetInt64(int64(score))
		}
		fmt.Println(i, ballot)
	}

	var elected [5]bool
	for round := 1; round <= 4; round++ {
		// Loop invariant: elected has round-1 indices as true.

		// Sum reweighted scores for each candidate.
		var totals [5]big.Rat
		for _, ballot := range ballots {
			for j, score := range ballot.scores {
				totals[j].Add(&totals[j], new(big.Rat).Quo(&score, &ballot.invWeight))
			}
		}
		fmt.Print("Totals:")
		for _, total := range totals {
			f, _ := total.Float64()
			fmt.Printf(" %.3f", f)
		}
		fmt.Println()

		// Find unelected candidate with highest score.
		best := -1
		for j, score := range totals {
			if elected[j] {
				continue
			}
			if best < 0 || score.Cmp(&totals[best]) > 0 {
				best = j
			}
		}

		// Add them as a winner.
		fmt.Println("Next winner:", best)
		elected[best] = true

		// Adjust weights.
		for i := range ballots {
			b := &ballots[i]
			b.invWeight.Add(&b.invWeight, new(big.Rat).Quo(&b.scores[best], new(big.Rat).SetInt64(9)))
		}
	}

	fmt.Println(elected)
}

type ballot struct {
	invWeight big.Rat    // inverse weight; 1 + electedScores / 9*elected
	scores    [5]big.Rat // raw scores for each candidate
}

func (b ballot) String() string {
	return fmt.Sprintf("{InvWeight: %v, Rocky: %v, Kitty: %v, Cassie: %v, Zoe: %v, Susana: %v}", b.invWeight.String(), b.scores[0].String(), b.scores[1].String(), b.scores[2].String(), b.scores[3].String(), b.scores[4].String())
}
