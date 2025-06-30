package main

import (
	"fmt"
	"strings"
)

func main() {
	PalindromicSubstring()

}

func PalindromicSubstring() {
	s := "babad"
	s2 := "abb"

	fmt.Printf("%v - %T \n", s[:], s[:])
	fmt.Printf("s2 %v - \n", s2[1:3])

	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		a := string(s[i])
		b := string(s[j])
		fmt.Printf("a=%v, b=%v \n", a, b)

		if a != b {
			fmt.Println("NOPT EQUAL!")
			break
		}

	}

	a := [3]rune{'b', 'a', 'b'}
	b := [3]rune{'b', 'a', 'b'}

	fmt.Println(a == b)

	/*

		"babad"
			b == d X

		"baba"
		"bab"
		"ba"a


		"abad"

	*/

}

func substringSearch() {
	//scenario := "dvdf"
	scenario := "pwwkew"
	var longest int
	var current int
	var substring string
	for _, r := range scenario {
		c := string(r)
		//fmt.Println(c)

		if strings.Contains(substring, c) {
			substringNew := ""
			for _, v := range substring {
				c2 := string(v)
				if c2 == c {
					substringNew = ""
					continue
				}
				substringNew += c2
			}

			fmt.Println(substring, substringNew)
			substring = substringNew
			current = len(substring)

		}

		substring += c
		current++
		if current > longest {
			longest = current
		}
	}

	fmt.Printf("longest=%v \n", longest)
}
