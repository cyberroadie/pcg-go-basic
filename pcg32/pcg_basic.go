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

/*
 * This code is derived from the minimal C implmenation which is in turn derived
 * from the full C implementation, which is in turn
 * derived from the canonical C++ PCG implementation. The C++ version
 * has many additional features and is preferable if you can use C++ in
 * your project.
 */
package pcg32

type RandomStruct struct {
	State uint64
	Inc   uint64
}

var PCG32_INITIALIZER = RandomStruct{0x853c49e6748fea9b, 0xda3e39cb94b95bdb}

//     Seed the rng.  Specified in two parts, state initializer and a
//     sequence selection constant (a.k.a. stream id)

func SeedGlobal(initstate uint64, initseq uint64) {
	Seed(&PCG32_INITIALIZER, initstate, initseq)
}

func Seed(rng*RandomStruct, initstate uint64, initseq uint64) {
	rng.State = 0
	rng.Inc = (initseq << 1) | 1
	Random(rng)
	rng.State += initstate
	Random(rng)
}

//     Generate a uniformly distributed 32-bit random number
func Random(rng*RandomStruct) uint32 {
	oldstate := rng.State
	rng.State = oldstate * 6364136223846793005 + rng.Inc
	xorshifted := uint32(((oldstate >> 18) ^ oldstate) >> 27)
	rot := uint32(oldstate >> 59)
	return (xorshifted >> rot) | (xorshifted << ((-rot) & 31))
}

func RandomGlobal() uint32 {
	return Random(&PCG32_INITIALIZER)
}

//     Generate a uniformly distributed number, r, where 0 <= r < bound
func BoundedRandom(rng*RandomStruct, bound uint32) uint32 {
	// To avoid bias, we need to make the range of the RNG a multiple of
	// bound, which we do by dropping output less than a threshold.
	// A naive scheme to calculate the threshold would be to do
	//
	//     uint32_t threshold = 0x100000000ull % bound;
	//
	// but 64-bit div/mod is slower than 32-bit div/mod (especially on
	// 32-bit platforms).  In essence, we do
	//
	//     uint32_t threshold = (0x100000000ull-bound) % bound;
	//
	// because this version will calculate the same modulus, but the LHS
	// value is less than 2^32.

	threshold := -bound % bound;

	// Uniformity guarantees that this loop will terminate.  In practice, it
	// should usually terminate quickly; on average (assuming all bounds are
	// equally likely), 82.25% of the time, we can expect it to require just
	// one iteration.  In the worst case, someone passes a bound of 2^31 + 1
	// (i.e., 2147483649), which invalidates almost 50% of the range.  In
	// practice, bounds are typically small and only a tiny amount of the range
	// is eliminated.
	for {
		r := Random(rng)
		if r >= threshold {
			return r % bound
		}
	}
}