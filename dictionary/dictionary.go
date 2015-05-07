/*
Package dictionary provides an interface to a dictionary with fancy features
such as complex letter frequency.

TODO: add complex letter frequency

Author: Justin Cook <jhcook@gmail.com>
*/

package dictionary

import (
	"github.com/jhcook/game_engine/dictionary/cache_sqlite"
	"github.com/jhcook/game_engine/util"
	"launchpad.net/xmlpath"
	"fmt"
	"bufio"
	"encoding/json"
	"errors"
	"time"
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
		} else {
			log.Println(line, "doesn't match")
		}
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


var ak string = ""
var url string = "https://dictionary.cambridge.org/api/v1/"

type DictEntry struct {
	Word        string
	Definition  string
}

/*
Returns spelling suggestions for a word, i.e. original and inflected forms
which resemble the word entered. Note that this will return suggestions even
for words spelled correctly!

Request Parameters
Accept [STRING] Header; Default: application/json application/json or
                                 application/xml
dictCode [IDSTRING] URL Dictionary code, which can be found from the
                    getDictionaries method
entrynumber [INT] Parameter Number of items to be shown.
q [STRING] Parameter The string to search for.

Response 200
{
  dictionaryCode [IDSTRING] Unique id of the dictionary dataset.
  searchTerm [STRING] The term that the user has entered.
  suggestions[STRING] The suggested new search term.
}
*/

func didYouMeanCambridge(wrd string) ([]interface{}, error) {
	var srch map[string]interface{}
	api := "dictionaries/american-english/search/didyoumean?q="

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url+api+wrd, nil)
	req.Header.Set("Host", "dictionary.cambridge.org")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("accessKey", ak)

	res, err := httpClient.Do(req)
	if err != nil {
		log.Println(fmt.Sprintf("%s.Do: %s", util.GetFuncName(), err))
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &srch); err != nil {
		log.Println(fmt.Sprintf("%s.Unmarshal: %s", util.GetFuncName(), err))
		return nil, err
	} else {
		res.Body.Close()
	}

	if wrds, ok := srch["suggestions"].([]interface{}); ok {
		return wrds, nil
	}
	return nil, errors.New(fmt.Sprintf("%s: unable to source words", util.GetFuncName()))
}

/*
Performs a search for a word or phrase in a particular dictionary.

currentPageIndex [INT] The index (offset) of the current page of results.
dictionaryCode [IDSTRING] Unique id of the dictionary dataset.
pageNumber [INT] The total number of pages of results found.
resultNumber [INT] The total number of results found.
results
[
  {
    entryId [IDSTRING] Unique id of the dictionary entry.
    entryLabel [STRING] A user-facing label which describes the entry.
    entryUrl [URL] The canonical URL of the entry on Cambridge Dictionaries Online.
  }
]
*/

func searchCambridge(wrd string) ([]interface{}, error) {
	var srch map[string]interface{}
	api := "dictionaries/american-english/search/?q="

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url+api+wrd, nil)
	req.Header.Set("Host", "dictionary.cambridge.org")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("accessKey", ak)

	res, err := httpClient.Do(req)
	if err != nil {
		log.Println(fmt.Sprintf("%s.Do: %s", util.GetFuncName(), err))
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &srch); err != nil {
		log.Println(fmt.Sprintf("%s.Unmarshal: %s", util.GetFuncName(), err))
		return nil, err
	} else {
		res.Body.Close()
	}

	if wrds, ok := srch["results"].([]interface{}); ok {
		return wrds, nil
	}
	return nil, errors.New(fmt.Sprintf("%s: unable to source words", util.GetFuncName()))
}

/*
Queries the API for `wrd` and receives on 200:

dictionaryCode [IDSTRING] Unique id of the dictionary dataset.
entryContent [STRING] The actual content of an entry.
entryId [IDSTRING] Unique id of the dictionary entry.
entryLabel [STRING] A user-facing label which describes the entry.
entryUrl [URL] The canonical URL of the entry on Cambridge Dictionaries Online.
format [STRING] Format of entry content returned.
topics
[
  {
    topicId [IDSTRING] Unique id of the individual topic.
    topicLabel [STRING] A user-facing label which describes the topic.
    topicParentId [IDSTRING] Unique id of the parent of the topic.
    topicUrl [STRING] The canonical URL of the topic on Cambridge Dictionaries Online.
  }
]

getDefinitionCambridge parses entryContent for
<sense-block><def-block><definition> and returns the string in <def>.
*/

func getDefinitionCambridge(wrd string) (string, error) {
	var msg map[string]interface{}
	api := "dictionaries/american-english/entries/"
	opts := "?format=xml"

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url+api+wrd+opts, nil)
	req.Header.Set("Host", "dictionary.cambridge.org")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("accessKey", ak)

	res, err := httpClient.Do(req)
	if err != nil {
		log.Println(fmt.Sprintf("%s.Do: %s", util.GetFuncName(), err))
		return "", err
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &msg); err != nil {
		log.Println(fmt.Sprintf("%s: %s", util.GetFuncName(), err))
		return "", err
	} else {
		res.Body.Close()
	}

	// A successful HTTP GET was performed, but see if word was found.
	content, isOk := msg["entryContent"].(string)

	// If entryContent does not exist, the entry does not exist.
	if !isOk {
		emptyStringErr := errors.New(fmt.Sprintf("%s: no entryContent", util.GetFuncName()))
		return "", emptyStringErr
	}

	// The entry exists so parse xml for definition.
	entryContent := strings.NewReader(content)
	path := xmlpath.MustCompile("/di/pos-block/sense-block/def-block[1]/definition/def")
	root, err := xmlpath.Parse(entryContent)
	if err != nil {
		return "", errors.New(fmt.Sprintf("%s.xmlpath.Parse: %s", util.GetFuncName(), err))
	}
	if value, ok := path.String(root); ok {
		return strings.Trim(value, ": "), nil
	} else {
		return "", errors.New(fmt.Sprintf("%s.xmlpath.String: value not found",
			util.GetFuncName()))
	}
}

/*
This is the function that performs the orchestration using the Cambridge
remote dictionary API.
*/
func (d *DictEntry) GetDefinition() (error) {
	defer func() (string, error) {
		if r := recover(); r != nil {
			return "", errors.New(fmt.Sprintf("%s: unable to source", util.GetFuncName()))
		}
		return "", errors.New(fmt.Sprintf("%s: unknown error", util.GetFuncName()))
	}()

	// Get dictionary entry
	var err error

	// Check to see if the word is in the cache
    log.Println(fmt.Sprintf("%s: checking for %s", util.GetFuncName(), d.Word))
	request <- d.Word
	stf := <- outChan
	if stf != nil { // && stf.Definition.Valid {
        log.Println(fmt.Sprintf("%v"), stf)
        if stf.Definition.Valid {
            log.Println(fmt.Sprintf("%s: %s is in cache", util.GetFuncName(), d.Word))
		    d.Definition = stf.Definition.String
		    return nil
        }
	}

	d.Definition, err = getDefinitionCambridge(d.Word)

	/* If successful: content received. Otherwise, perform inflection and
	 * search for any matching words.
	 */
	if err != nil {
		log.Println(err)
		if wrds, err2 := didYouMeanCambridge(d.Word); err2 != nil {
			log.Println(err2)
			return err2
		} else {
			d.Word = wrds[0].(string)
			wrds2, err3 := searchCambridge(d.Word)
			if err3 != nil {
				log.Println("err3: %v", err3)
				return errors.New(fmt.Sprintf("%s.searchCambridge: unable to source words", util.GetFuncName()))
			}
			d.Word = wrds2[0].(map[string]interface{})["entryId"].(string)
			var err4 error
			d.Definition, err4 = getDefinitionCambridge(d.Word)
			if err4 != nil {
				log.Println("err4: %s", err4)
				return errors.New(fmt.Sprintf("%s.getDefinitionCambridge: unable to source word", util.GetFuncName()))
			}
		}
	}
	dent := []string{d.Word, d.Definition}
	inpChan <- dent
	return nil
}


