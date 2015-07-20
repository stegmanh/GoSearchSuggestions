package trie

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var plainWord = regexp.MustCompile(`(^[a-zA-Z_ ]*$)`)
var alreadyAdded = make(map[string]bool)

type Trie struct {
	root *trieNode
}

type trieNode struct {
	value     string
	container *[]string
	letters   *[27]*trieNode
}

func (t *Trie) Initialize() {
	temp := [27]*trieNode{}
	t.root = &trieNode{value: "", letters: &temp}
}

func (t *Trie) BuildTrie(scanner *bufio.Scanner) {
	count := 0
	for scanner.Scan() {
		s := strings.ToLower(scanner.Text())
		matched := plainWord.MatchString(s)
		if matched {
			s = strings.Replace(s, "_", " ", -1)
			t.Add(s)
			count++
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error adding to trie: ", err)
	}
	fmt.Println("Build Trie")
}

func (t *Trie) Add(s string) {
	current := t.root
	idx := 0
	for idx < len(s) {
		pos := s[idx] - 97
		//Catch the space's
		if pos < 0 || pos > 25 {
			pos = 26
		}
		idx++
		if current.container != nil {
			tempContainer := *current.container
			tempContainer = append(tempContainer, s)
			current.container = &tempContainer
			if len(*current.container) >= 20 {
				current.letters = &[27]*trieNode{}
				for _, str := range *current.container {
					// Check for space or whatever char code 191 + 97 is
					insertPos := str[idx-1] - 97
					if insertPos < 0 || insertPos > 25 {
						insertPos = 26
					}
					if current.letters[insertPos] == nil {
						temp := make([]string, 0)
						current.letters[insertPos] = &trieNode{value: "", container: &temp}
					}
					if idx == len(str) {
						current.letters[insertPos].value = str
						continue
					}
					tCurrent := current.letters[insertPos]
					tempContainer := *tCurrent.container
					tempContainer = append(tempContainer, str)
					tCurrent.container = &tempContainer
				}
				current.container = nil
			}
			return
		}
		if current.letters[pos] == nil {
			temp := make([]string, 0)
			current.letters[pos] = &trieNode{value: "", container: &temp}
			current = current.letters[pos]
			if idx == len(s) {
				current.value = s
			} else {
				tempContainer := *current.container
				tempContainer = append(tempContainer, s)
				current.container = &tempContainer
			}
		} else {
			current = current.letters[pos]
		}
	}
}

func (t *Trie) Find(searchString string) []string {
	searchString = strings.ToLower(searchString)
	if !plainWord.MatchString(searchString) {
		return make([]string, 0)
	}
	current := t.root
	for idx := 0; idx < len(searchString); idx++ {
		pos := searchString[idx] - 97
		if pos < 0 || pos > 26 {
			pos = 26
		}
		if current.letters != nil {
			if current.letters[pos] == nil {
				return make([]string, 0)
			}
			current = current.letters[pos]
		} else {
			//Meh this should change, doesnt look good
			break
		}
	}
	if current.container != nil {
		return filter(*current.container, searchString)
	}
	toReturn := make([]string, 0)
	if current.value != "" {
		toReturn = append(toReturn, current.value)
	}
	return findHelper(current, &toReturn)
}

func findHelper(start *trieNode, toReturn *[]string) []string {
	if start.container != nil {
		for _, val := range *start.container {
			*toReturn = append(*toReturn, val)
		}
		//If letters < 9 -- Add later
		return *toReturn
	}
	for _, pStart := range *start.letters {
		if len(*toReturn) > 9 {
			return *toReturn
		}
		if pStart != nil {
			if pStart.value != "" {
				*toReturn = append(*toReturn, pStart.value)
			}
			findHelper(pStart, toReturn)
		}
	}
	return *toReturn
}

func filter(s []string, start string) (toReturn []string) {
	for _, word := range s {
		if strings.HasPrefix(word, start) {
			toReturn = append(toReturn, word)
		}
	}
	return
}
