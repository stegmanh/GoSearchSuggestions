package trie

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strings"
)

var plainWord = regexp.MustCompile(`(^[a-zA-Z]*$)`)

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
		if !matched {
			continue
		}
		t.Add(s)
		count++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done")
}

func (t *Trie) Add(s string) {
	current := t.root
	idx := 0
	for idx < len(s) {
		pos := s[idx] - 97
		idx++
		if current.container != nil {
			tempContainer := *current.container
			tempContainer = append(tempContainer, s)
			current.container = &tempContainer
			if len(*current.container) >= 20 {
				current.letters = &[27]*trieNode{}
				for _, str := range *current.container {
					// fmt.Printf("%v, %#v\n", idx, str)
					insertPos := str[idx-1] - 97
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
	if !plainWord.MatchString(searchString) {
		fmt.Println("Not valid")
		return make([]string, 0)
	}
	current := t.root
	for idx := 0; idx < len(searchString); idx++ {
		pos := searchString[idx] - 97
		if current.container != nil {
			fmt.Println("Returning container...")
			return *current.container
		}
		if current.letters[pos] == nil {
			fmt.Println("Returning this weird...")
			return make([]string, 0)
		}
		current = current.letters[pos]
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
