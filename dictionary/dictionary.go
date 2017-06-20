/*
Package dictionary provides an interface to a dictionary with fancy features
such as complex letter frequency.

TODO: add complex letter frequency

Author: Justin Cook <jhcook@gmail.com>
*/

package dictionary

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jhcook/game_engine/dictionary/cache_sqlite"
	"github.com/jhcook/game_engine/util"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"time"
)

type Dictionary struct {
	Word  string
	Ci    uint32
	Words []string
}

var inpChan = make(chan []string)
var outChan = make(chan *cache_sqlite.WordDefinition)
var request = make(chan string)

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln("%s: %v\n", util.GetFuncName(), err)
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if match, _ := regexp.MatchString("^[A-Za-z]{5,}$", line); match {
			lines = append(lines, scanner.Text())
		} /*else {
			log.Println(line, "doesn't match")
		} */
	}
	return lines, scanner.Err()
}

func NewDictionary(path string) *Dictionary {

	// Instantiate cache
	go cache_sqlite.DefinitionWriter(inpChan, request, outChan)
	// Need time to create if db does not exist
	time.Sleep(1 * time.Second)
	go cache_sqlite.DefinitionReader(request, outChan)

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

var ak string = os.Getenv("__AK__")
var aid string = os.Getenv("__AID__")
var url string = "https://od-api.oxforddictionaries.com:443/api/v1"

type DictEntry struct {
	Word       string
	Definition string
}

/*
Queries the API for `wrd` and receives on 200:

dictionaryCode [IDSTRING] Unique id of the dictionary dataset.
entryContent [STRING] The actual content of an entry.
entryId [IDSTRING] Unique id of the dictionary entry.
entryLabel [STRING] A user-facing label which describes the entry.
entryUrl [URL] The canonical URL of the entry on Oxford Dictionaries Online.
format [STRING] Format of entry content returned.
topics
[
  {
    topicId [IDSTRING] Unique id of the individual topic.
    topicLabel [STRING] A user-facing label which describes the topic.
    topicParentId [IDSTRING] Unique id of the parent of the topic.
    topicUrl [STRING] The canonical URL of the topic on Oxford Dictionaries Online.
  }
]

getDefinitionOxford parses entryContent for
<sense-block><def-block><definition> and returns the string in <def>.
*/

func getDefinitionOxford(wrd string) (string, error) {
	var msg map[string]interface{}
	api := "/entries/en/"
	opts := "/definitions"

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url+api+wrd+opts, nil)
	//req.Header.Set("Host", "dictionary.oxford.org")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("app_id", aid)
	req.Header.Set("app_key", ak)

	res, err := httpClient.Do(req)
	if err != nil {
		log.Println(fmt.Sprintf("%s.Do: %s", util.GetFuncName(), err))
		return "", err
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &msg); err != nil {
		log.Println(fmt.Sprintf("%s: %s", util.GetFuncName(), err))
		fmt.Println(string(body))
		return "", err
	} else {
		res.Body.Close()
	}

	// The entry exists so parse json for definition -- yes this is ugly
	// results
	s := reflect.ValueOf(msg["results"])
	results := make([]interface{}, s.Len())
	for i := range msg["results"].([]interface{}) {
		results[i] = s.Index(i).Interface()
	}

	s = reflect.ValueOf(results[0])
	firstResult := make(map[string]interface{}, s.Len())
	for k, v := range results[0].(map[string]interface{}) {
		firstResult[k] = v
	}

	// lexicalEntries
	s = reflect.ValueOf(firstResult["lexicalEntries"])
	lexicalEntries := make([]interface{}, s.Len())
	for i := range firstResult["lexicalEntries"].([]interface{}) {
		lexicalEntries[i] = s.Index(i).Interface()
	}

	s = reflect.ValueOf(lexicalEntries[0])
	firstLexicalEntry := make(map[string]interface{}, s.Len())
	for k, v := range lexicalEntries[0].(map[string]interface{}) {
		firstLexicalEntry[k] = v
	}

	// entries
	s = reflect.ValueOf(firstLexicalEntry["entries"])
	entries := make([]interface{}, s.Len())
	for i := range firstLexicalEntry["entries"].([]interface{}) {
		entries[i] = s.Index(i).Interface()
	}

	s = reflect.ValueOf(entries[0])
	firstEntry := make(map[string]interface{}, s.Len())
	for k, v := range entries[0].(map[string]interface{}) {
		firstEntry[k] = v
	}

	// senses
	s = reflect.ValueOf(firstEntry["senses"])
	senses := make([]interface{}, s.Len())
	for i := range firstEntry["senses"].([]interface{}) {
		senses[i] = s.Index(i).Interface()
	}

	s = reflect.ValueOf(senses[0])
	firstSense := make(map[string]interface{}, s.Len())
	for k, v := range senses[0].(map[string]interface{}) {
		firstSense[k] = v
	}

	// definitions
	s = reflect.ValueOf(firstSense["definitions"])
	definitions := make([]interface{}, s.Len())
	for i := range firstSense["definitions"].([]interface{}) {
		definitions[i] = s.Index(i).Interface()
	}

	return definitions[0].(string), err
}

// Check to see if the word is in the cache
func (d *DictEntry) checkCache() error {
	log.Printf("%s: checking for %s", util.GetFuncName(), d.Word)
	request <- d.Word
	stf := <-outChan
	if stf != nil { // && stf.Definition.Valid {
		log.Println(fmt.Sprintf("%v"), stf)
		if stf.Definition.Valid {
			log.Println(fmt.Sprintf("%s: %s is in cache", util.GetFuncName(), d.Word))
			d.Definition = stf.Definition.String
			return nil
		}
	}
	return errors.New(fmt.Sprintf("%s: %s not in cache", util.GetFuncName(), d.Word))
}

/*
This is the function that performs the orchestration using the Oxford
remote dictionary API.
*/
func (d *DictEntry) GetDefinition() error {
	defer func() (string, error) {
		if r := recover(); r != nil {
			return "", errors.New(fmt.Sprintf("%s: unable to source", util.GetFuncName()))
		}
		return "", errors.New(fmt.Sprintf("%s: unknown error", util.GetFuncName()))
	}()

	var err error
	// Check to see if the word is in the cache
	if err = d.checkCache(); err != nil {
		log.Println(err)
		d.Definition, err = getDefinitionOxford(d.Word)
	} else {
		return nil
	}

	dent := []string{d.Word, d.Definition}
	inpChan <- dent
	return nil
}
