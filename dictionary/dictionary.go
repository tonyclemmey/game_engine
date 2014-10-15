/*
Package dictionary provides an interface to a dictionary with fancy features
such as complex letter frequency.

Author: Justin Cook <jhcook@gmail.com>
*/

package dictionary

import (
	"bufio"
	"code.google.com/p/go.net/html"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Dictionary struct {
	Word  string
	Ci    uint32
	Words []string
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
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if match, _ := regexp.MatchString("^[A-Za-z]{5,}$", line); match {
			lines = append(lines, scanner.Text())
		} else {
			log.Println(line, "doesn't match")
		}
	}
	return lines, scanner.Err()
}

func NewDictionary(path string) *Dictionary {
	dict := new(Dictionary)
	dict.Ci = 1000
	if len(path) == 0 {
		path = "/usr/share/dict/words"
	}
	dict.Words, _ = readLines(path)
	return dict
}

func (d *Dictionary) NextWord() string {
	d.Ci++
	d.Word = d.Words[d.Ci]
	return d.Word
}

func SearchMacmillan(wrd string) (string, error) {
	url := "https://www.macmillandictionary.com/api/v1/"
	ak := "jK67wm71vm0PjdxolDZtKMMMIijzaSuxXslJPWcP50Vq87RWXW0SkLNS1sLRDBc4"
	var srch map[string]interface{}

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url+"dictionaries/american/search/?q="+wrd, nil)
	req.Header.Set("Host", "www.macmillandictionary.com")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("accessKey", ak)

	res, err := httpClient.Do(req)
	if err != nil {
		log.Println("dictionary.SearchMacmillan.Do:", err)
		return "", err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &srch); err != nil {
		log.Println("dictionary.SearchMacmillan.Unmarshal:", err)
		return "", err
	} else {
		res.Body.Close()
	}

	num, ok := srch["resultNumber"]

	var num1 int
	switch num.(type) {
	default:
		num1 = int(num.(float64)) //Assuming float64 (not portable)
	}

	if num1 < 1 {
		return "", errors.New("dictionary.SearchMacmillann: word not found in search")
	} else if !ok {
		return "", errors.New("dictionary.SearchMacmillan: element not found in search")
	}

	entry_map := srch["results"].([]interface{})[0].(map[string]interface{})
	return entry_map["entryId"].(string), nil
}

func GetMacmillan(wrd string) (map[string]interface{}, error) {
	url := "https://www.macmillandictionary.com/api/v1/"
	ak := "jK67wm71vm0PjdxolDZtKMMMIijzaSuxXslJPWcP50Vq87RWXW0SkLNS1sLRDBc4"
	var msg map[string]interface{}

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url+"dictionaries/american/entries/"+wrd, nil)
	req.Header.Set("Host", "www.macmillandictionary.com")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("accessKey", ak)

	res, err := httpClient.Do(req)
	if err != nil {
		log.Println("dictionary.GetMacmillan.Do:", err)
		return nil, err
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &msg); err != nil {
		log.Println("dictionary.GetMacmillan.Unmarshal:", err)
		return nil, err
	} else {
		res.Body.Close()
	}
	return msg, nil
}

func GetDefinition(wrd string) (string, error) {
	log.Println("dictionary.GetDefinition:", wrd)
	defer func() (string, error) {
		if r := recover(); r != nil {
			return "", errors.New("GetDefinition: unable to source")
		}
		return "", errors.New("GetDefinition: unknown error")
	}()

	// Search for word inflections
	new_word, err := SearchMacmillan(wrd)
	if err != nil {
		return "", err
	}

	// Get dictionary entry
	msg, err := GetMacmillan(new_word)
	if err != nil {
		return "", err
	}

	var f func(*html.Node) (string, bool)
	f = func(n *html.Node) (string, bool) {
		if n.Type == html.ElementNode && n.Data == "span" {
			for _, a := range n.Attr {
				if match, _ := regexp.MatchString(`DEFINITION`, a.Val); match {
					return n.FirstChild.Data, true
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if b, ok := f(c); ok {
				return b, true
			}
		}
		return "", false
	}

	doc, err := html.Parse(strings.NewReader(msg["entryContent"].(string)))
	if err != nil {
		log.Println("dictionary.GetDefinition.Parse:", err)
		return "", err
	}
	if boom, ok := f(doc); ok {
		re, _ := regexp.Compile(`[\w,\ \n]+`)
		res := re.FindAllStringSubmatch(boom, -1)
		return res[0][0], nil
	} else {
		return "", errors.New("dictionary.GetDefinition.regexp: failed to find matching string")
	}
}
