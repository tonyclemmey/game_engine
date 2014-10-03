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

func GetDefinition(wrd string) (string, error) {
    url := "https://www.macmillandictionary.com/api/v1/"
    ak := "jK67wm71vm0PjdxolDZtKMMMIijzaSuxXslJPWcP50Vq87RWXW0SkLNS1sLRDBc4"

    log.Println("dictionary.GetDefinition:", wrd)
    defer func() (string, error) {
        if r := recover(); r != nil {
            return "", errors.New("GetDefinition: unable to source")
        }
        return "", errors.New("GetDefinition: uknown error")
    }()

    //var msg []map[string]interface{}
    var msg map[string]interface{}

    httpClient := &http.Client{}

    req, err := http.NewRequest("GET", url+"dictionaries/american/entries/"+wrd, nil)
    req.Header.Set("Host", "www.macmillandictionary.com")
    req.Header.Set("Accept", "application/json")
    req.Header.Set("accessKey", ak)

    res, err := httpClient.Do(req)
    if err != nil {
        log.Println("dictionary.GetDefinition.Do:", err)
        return "", err
    }

    defer res.Body.Close()

    body, _ := ioutil.ReadAll(res.Body)

    if err := json.Unmarshal(body, &msg); err != nil {
        log.Println("dictionary.GetDefinition.Unmarshal:", err)
        return "", err
    } else {
        res.Body.Close()
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
