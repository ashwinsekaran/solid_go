package main

import (
	"fmt"
	"sort"
	"strings"
)

// ═══════════════════════════════════════════════════════════════════
// APPROACH 1 — Fixed [26]int Array (Best approach)
// Time: O(N×M)  Space: O(1)
// Most idiomatic Go for lowercase letter problems
// ═══════════════════════════════════════════════════════════════════
func solveFixedArray(N int, words []string) string {
	bestWord := ""
	bestCount := -1

	for _, word := range words {
		// [26]int — index 0='a', index 25='z', all start at 0
		var freq [26]int
		for _, ch := range word {
			freq[ch-'a']++ // 'b'-'a'=1, so freq[1]++ for 'b'
		}

		uniqueCount := 0
		for _, c := range freq {
			if c == 1 { // exactly once = unique letter
				uniqueCount++
			}
		}

		if uniqueCount > bestCount ||
			(uniqueCount == bestCount && word < bestWord) {
			bestCount = uniqueCount
			bestWord = word
		}
	}
	return bestWord
}

// ═══════════════════════════════════════════════════════════════════
// APPROACH 2 — HashMap map[rune]int
// Time: O(N×M)  Space: O(M)
// More flexible — works for any character set (unicode, special chars)
// ═══════════════════════════════════════════════════════════════════
func solveHashMap(N int, words []string) string {
	bestWord := ""
	bestCount := -1

	for _, word := range words {
		// map grows dynamically — heap allocation per word
		freq := make(map[rune]int)
		for _, ch := range word {
			freq[ch]++ // no index math needed
		}

		uniqueCount := 0
		for _, count := range freq {
			if count == 1 {
				uniqueCount++
			}
		}

		if uniqueCount > bestCount ||
			(uniqueCount == bestCount && word < bestWord) {
			bestCount = uniqueCount
			bestWord = word
		}
	}
	return bestWord
}

// ═══════════════════════════════════════════════════════════════════
// APPROACH 3 — Bitmask (Two uint32 integers)
// Time: O(N×M)  Space: O(1)
// Fastest in absolute time — uses CPU-level bitwise operations
// seen    = bit set if letter appeared >= 1 time
// duplicate = bit set if letter appeared >= 2 times
// unique  = seen AND NOT duplicate
// ═══════════════════════════════════════════════════════════════════
func solveBitmask(N int, words []string) string {
	bestWord := ""
	bestCount := -1

	for _, word := range words {
		var seen uint32      // tracks: have we seen this letter?
		var duplicate uint32 // tracks: have we seen this letter TWICE?

		for _, ch := range word {
			bit := uint32(1) << (ch - 'a') // e.g. 'a'=bit0, 'b'=bit1
			if seen&bit != 0 {
				// already seen before → now it's a duplicate
				duplicate |= bit
			}
			seen |= bit // mark as seen
		}

		// unique letters = seen but NOT in duplicate
		// &^ is Go's bit clear operator (AND NOT)
		unique := seen &^ duplicate

		// count set bits using Brian Kernighan's algorithm
		// each iteration clears the lowest set bit
		uniqueCount := 0
		n := unique
		for n != 0 {
			n &= n - 1 // clears lowest set bit
			uniqueCount++
		}

		if uniqueCount > bestCount ||
			(uniqueCount == bestCount && word < bestWord) {
			bestCount = uniqueCount
			bestWord = word
		}
	}
	return bestWord
}

// ═══════════════════════════════════════════════════════════════════
// APPROACH 4 — Sort Characters
// Time: O(N×M log M)  Space: O(M)
// After sorting, duplicates are adjacent — easy to spot
// Use when character set is unknown (unicode, multilingual)
// ═══════════════════════════════════════════════════════════════════
func solveSortChars(N int, words []string) string {
	bestWord := ""
	bestCount := -1

	for _, word := range words {
		// convert string to rune slice so we can sort
		chars := []rune(word)
		sort.Slice(chars, func(i, j int) bool {
			return chars[i] < chars[j]
		})

		// after sorting: "once" → "ceno"
		// walk through groups of identical chars
		uniqueCount := 0
		i := 0
		for i < len(chars) {
			j := i + 1
			// skip all identical consecutive chars
			for j < len(chars) && chars[j] == chars[i] {
				j++
			}
			// if group size is 1 → unique letter
			if j-i == 1 {
				uniqueCount++
			}
			i = j // move to next group
		}

		if uniqueCount > bestCount ||
			(uniqueCount == bestCount && word < bestWord) {
			bestCount = uniqueCount
			bestWord = word
		}
	}
	return bestWord
}

// ═══════════════════════════════════════════════════════════════════
// APPROACH 5 — strings.Count (Most readable, worst performance)
// Time: O(N×M²)  Space: O(M)
// strings.Count scans the ENTIRE word for each character
// Good for readability, never use in competitive programming
// ═══════════════════════════════════════════════════════════════════
func solveFunctional(N int, words []string) string {
	bestWord := ""
	bestCount := -1

	for _, word := range words {
		uniqueCount := 0
		seen := make(map[rune]bool) // avoid rechecking same letter

		for _, ch := range word {
			if seen[ch] {
				continue // already counted this letter
			}
			seen[ch] = true
			// strings.Count scans entire word → O(M) per letter
			if strings.Count(word, string(ch)) == 1 {
				uniqueCount++
			}
		}

		if uniqueCount > bestCount ||
			(uniqueCount == bestCount && word < bestWord) {
			bestCount = uniqueCount
			bestWord = word
		}
	}
	return bestWord
}

// ═══════════════════════════════════════════════════════════════════
// TEST RUNNER
// ═══════════════════════════════════════════════════════════════════
type approach struct {
	name string
	fn   func(int, []string) string
}

type testCase struct {
	name  string
	words []string
	want  string
}

func runTests(approaches []approach, tests []testCase) {
	totalApproaches := len(approaches)
	totalTests := len(tests)
	overallPass := 0

	for _, ap := range approaches {
		fmt.Printf("\n┌─────────────────────────────────────────────────────┐\n")
		fmt.Printf("│  %-51s│\n", ap.name)
		fmt.Printf("└─────────────────────────────────────────────────────┘\n")

		passed := 0
		for _, tc := range tests {
			got := ap.fn(len(tc.words), tc.words)
			if got == tc.want {
				passed++
				fmt.Printf("  ✅ PASS │ %-35s │ got: %q\n", tc.name, got)
			} else {
				fmt.Printf("  ❌ FAIL │ %-35s │ got: %q  want: %q\n", tc.name, got, tc.want)
			}
		}
		fmt.Printf("  ─────────────────────────────────────────────────────\n")
		fmt.Printf("  Score: %d/%d\n", passed, totalTests)
		if passed == totalTests {
			overallPass++
		}
	}

	fmt.Printf("\n╔═════════════════════════════════════════════════════╗\n")
	fmt.Printf("║  SUMMARY: %d/%d approaches passed all test cases     ║\n", overallPass, totalApproaches)
	fmt.Printf("╚═════════════════════════════════════════════════════╝\n")
}

func main() {
	approaches := []approach{
		{"APPROACH 1 — Fixed [26]int Array  | O(N×M) time O(1) space", solveFixedArray},
		{"APPROACH 2 — HashMap map[rune]int | O(N×M) time O(M) space", solveHashMap},
		{"APPROACH 3 — Bitmask uint32       | O(N×M) time O(1) space", solveBitmask},
		{"APPROACH 4 — Sort Characters      | O(N×M logM) time O(M) space", solveSortChars},
		{"APPROACH 5 — strings.Count        | O(N×M²) time O(M) space", solveFunctional},
	}

	tests := []testCase{
		// ── from problem statement ──────────────────────────────
		{
			name:  "Sample1: once/three/sorry",
			words: []string{"once", "three", "sorry"},
			want:  "once",
			// once=4 unique, three=3, sorry=3 → once wins
		},
		{
			name:  "Sample2: trap/unhorizontal/gaet",
			words: []string{"trap", "unhorizontal", "gaet"},
			want:  "unhorizontal",
			// unhorizontal: u,n,h,r,i,z,o,t,a,l = 10 unique (o appears once)
		},
		{
			name:  "Sample3: single word gowan",
			words: []string{"gowan"},
			want:  "gowan",
			// all 5 letters unique
		},
		{
			name:  "Sample4: boarding/kevin/equalizers/sphaeriidae",
			words: []string{"boarding", "kevin", "equalizers", "sphaeriidae"},
			want:  "boarding",
			// boarding: b,o,a,r,d,i,n,g = 8 unique
		},

		// ── tie-breaking (lexicographic) ────────────────────────
		{
			name:  "Tie: abc vs abd → abc wins",
			words: []string{"abd", "abc"},
			want:  "abc",
			// both have 3 unique, abc < abd alphabetically
		},
		{
			name:  "Tie: apple vs apply → same unique count",
			words: []string{"apply", "apple"},
			want:  "apple",
			// apple: a=1,p=2,l=1,e=1 → 3 unique
			// apply: a=1,p=2,l=1,y=1 → 3 unique → apple < apply
		},
		{
			name:  "Tie: xyz vs abc → abc wins",
			words: []string{"xyz", "abc"},
			want:  "abc",
			// both 3 unique, abc < xyz
		},

		// ── edge cases ──────────────────────────────────────────
		{
			name:  "All duplicates: aabb vs xyz",
			words: []string{"aabb", "xyz"},
			want:  "xyz",
			// aabb: a=2,b=2 → 0 unique; xyz: all unique → xyz wins
		},
		{
			name:  "Single letter word: a",
			words: []string{"a"},
			want:  "a",
			// 'a' appears once → 1 unique letter
		},
		{
			name:  "All same letters: aaaa",
			words: []string{"aaaa", "b"},
			want:  "b",
			// aaaa: 0 unique; b: 1 unique → b wins
		},
		{
			name:  "Long word charmfully",
			words: []string{"charmfully"},
			want:  "charmfully",
			// c,h,a,r,m,f,u,l(×2),y → l is duplicate → 9 unique
		},
		{
			name:  "All words same unique count → first alphabetically",
			words: []string{"dog", "cat", "bat"},
			want:  "bat",
			// dog=3, cat=3, bat=3 → bat < cat < dog → bat wins
		},
		{
			name:  "One word all unique: abcdef",
			words: []string{"abcdef", "aabbcc"},
			want:  "abcdef",
			// abcdef=6 unique, aabbcc=0 unique → abcdef wins
		},
		{
			name:  "Repeated full word pattern: abab",
			words: []string{"abab", "cd"},
			want:  "cd",
			// abab: a=2,b=2 → 0 unique; cd: c=1,d=1 → 2 unique → cd wins
		},
	}

	fmt.Println("╔═════════════════════════════════════════════════════╗")
	fmt.Println("║     COUNT THE LETTERS — All Approaches Comparison   ║")
	fmt.Println("╚═════════════════════════════════════════════════════╝")
	fmt.Printf("\nRunning %d test cases across %d approaches...\n", len(tests), len(approaches))

	runTests(approaches, tests)

	// ── Complexity summary ──────────────────────────────────────
	fmt.Println("\n╔═════════════════════════════════════════════════════════════╗")
	fmt.Println("║                  COMPLEXITY COMPARISON                      ║")
	fmt.Println("╠══════════════════════════╦═════════════╦════════════════════╣")
	fmt.Println("║ Approach                 ║ Time        ║ Space              ║")
	fmt.Println("╠══════════════════════════╬═════════════╬════════════════════╣")
	fmt.Println("║ 1. Fixed [26]int Array   ║ O(N×M)      ║ O(1) ← best       ║")
	fmt.Println("║ 2. HashMap               ║ O(N×M)      ║ O(M)               ║")
	fmt.Println("║ 3. Bitmask uint32        ║ O(N×M)      ║ O(1) ← fastest    ║")
	fmt.Println("║ 4. Sort Characters       ║ O(N×M logM) ║ O(M)               ║")
	fmt.Println("║ 5. strings.Count         ║ O(N×M²)     ║ O(M)               ║")
	fmt.Println("╚══════════════════════════╩═════════════╩════════════════════╝")
	fmt.Println("\nN = number of words, M = average word length (max 100 per constraints)")
	fmt.Println("Recommended: Approach 1 for interviews, Approach 3 for max speed")
}
