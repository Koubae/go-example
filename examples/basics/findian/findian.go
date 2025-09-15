/*
Write a program which prompts the user to enter a string.
The program searches through the entered string
for the characters ‘i’, ‘a’, and ‘n’.

The program should print “Found!” if the entered string starts with the character ‘i’,
ends with the character ‘n’, and contains the character ‘a’.
The program should print “Not Found!” otherwise.
The program should not be case-sensitive, so it does not matter if the characters are upper-case or lower-case.

Examples:
The program should print “Found!” for the following example entered strings,
“ian”, “Ian”, “iuiygaygn”, “I d skd a efju N”.

The program should print “Not Found!” for the following strings,
“ihhhhhn”, “ina”, “xian”.
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const StartWith = 'i'
const EndWith = 'n'
const Contains = 'a'
const FoundMessage = "Found!"
const NotFoundMessage = "Not Found!"

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Get user input
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Normalize input by set it to lower case
	input = strings.ToLower(input)

	// Validate input format
	startsWithI := strings.HasPrefix(input, string(StartWith))
	endsWithI := strings.HasSuffix(input, string(EndWith))
	containsA := strings.Contains(input, string(Contains))

	// Process result
	result := NotFoundMessage
	if startsWithI && endsWithI && containsA {
		result = FoundMessage
	}

	// output result
	fmt.Println(result)

}
