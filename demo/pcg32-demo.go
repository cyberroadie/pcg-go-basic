/*
 * PCG Random Number Generation for Go.
 *
 * Copyright 2015 Olivier Van Acker <cyberroadie@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * For additional information about the PCG random number generation scheme,
 * including its license and other licensing options, visit
 *
 *       http://www.pcg-random.org
 */
package main


import (
	"flag"
	"fmt"
	"math/rand"
	"github.com/cyberroadie/pcg-go-basic/pcg32"
	"unsafe"
)

func main() {
	rounds := flag.Uint64("rounds", 5, "how many rounds")
	nondeterministicSeed := flag.Bool("r", false, "nondeterministic seed")
	global := flag.Bool("global", false, "global rng")

	var rng pcg32.RandomStruct

	if *global {
		rng = pcg32.PCG32_INITIALIZER
	}

	if *nondeterministicSeed {
		pcg32.Seed(&rng, uint64(rand.Uint32())<<32 + uint64(rand.Uint32()), *rounds) // Go 1 doesn't have rand.Uint64()
	} else {
		pcg32.Seed(&rng, uint64(42), uint64(54))
	}

	fmt.Printf("pcg32_random_r:\n" +
	"      -  result:      32-bit unsigned int (uint32)\n" +
	"      -  period:      2^64   (* 2^63 streams)\n" +
	"      -  state type:  pcg32.RandomStruct (%x bytes)\n" +
	"      -  output func: XSH-RR\n" +
	"\n", unsafe.Sizeof(rng))

	for round := uint64(1); round <= *rounds; round++ {
		fmt.Printf("Round %d:\n", round)

		/* Make some 32-bit numbers */
		fmt.Printf(" 32bit:")
		for i := 0; i < 6; i++ {
			fmt.Printf(" 0x%08x", pcg32.Random(&rng))
		}
		fmt.Println()

		/* Toss some coins */
		fmt.Printf(" Coins:")
		for i := 0; i < 65; i++ {
			switch pcg32.BoundedRandom(&rng, 2) {
			case 0:
				fmt.Print("H")
			case 1:
				fmt.Print("T")
			}
		}
		fmt.Println()

		/* Deal some cards */
		const (
			SUITS = uint32(4)
			NUMBERS = uint32(13)
			CARDS = uint32(52)
		)

		cards := [CARDS]uint32{}

		for i := uint32(0); i < CARDS; i++ {
			cards[i] = i
		}

		for i := CARDS; i > 1; i-- {
			chosen := pcg32.BoundedRandom(&rng, i)
			card := cards[chosen]
			cards[chosen] = cards[i - 1]
			cards[i - 1] = card
		}

		fmt.Printf("  Cards:")

		number := [NUMBERS]rune{'A', '2', '3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'Q', 'K'}
		suit := [SUITS]rune{'h', 'c', 'd', 's'}

		for i := uint32(0); i < CARDS; i++ {
			fmt.Printf(" %c%c", number[cards[i] / SUITS], suit[cards[i] % SUITS])
			if (i + 1) % 22 == 0 {
				fmt.Println("\t")
			}
		}
		fmt.Println()
		fmt.Println()
	}

}
