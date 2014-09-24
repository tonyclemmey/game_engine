/*
Package dictionary provides an interface to a dictionary with fancy features
such as complex letter frequency.

Author: Justin Cook <jhcook@gmail.com>
*/

package dictionary

import (
	"bufio"
	"log"
	"os"
)

type dictionary struct {
	word  string
	ci    uint32
	words []string
}

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

func NewDictionary(path string) *dictionary {
	dict := new(dictionary)
	dict.ci = 1000
	if len(path) == 0 {
		path = "/usr/share/dict/words"
	}
	dict.words, _ = readLines(path)
	return dict
}

func (d *dictionary) NextWord() string {
	d.ci++
	d.word = d.words[d.ci]
	return d.word
}
