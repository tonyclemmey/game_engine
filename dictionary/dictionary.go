/*
Package dictionary provides an interface to a dictionary with fancy features
such as complex letter frequency.

Author: Justin Cook <jhcook@gmail.com>
 */

package dictionary

import (
//	"fmt"
	"log"
	"os"
	"bufio"
)

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln("dictionary.readLines:", err)
		return nil, err
	}
	defer file.Close()
	log.Println("dictionary.readLines open:", path)
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
