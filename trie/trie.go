package trie

import (
	"bufio"
	"log"
	"regexp"
	"strings"
)

var plainWord = regexp.MustCompile(`(^[a-zA-Z]*$)`)

type Trie struct {
	root *trieNode
}

type trieNode struct {
	value   string
	letters [27]*trieNode
}

func (t *Trie) Initialize() {
	t.root = &trieNode{value: "", letters: [27]*trieNode{}}
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
}

func (t *Trie) Add(s string) {
	current := t.root
	idx := 0
	for idx < len(s) {
		pos := s[idx] - 97
		if current.letters[pos] == nil {
			current.letters[pos] = &trieNode{value: "", letters: [27]*trieNode{}}
		}
		current = current.letters[pos]
		idx++
		if idx == len(s) {
			current.value = s
		}
	}
}

func (t *Trie) Find(sub string) []string {
	if !plainWord.MatchString(sub) {
		return make([]string, 0)
	}
	current := t.root
	for idx := 0; idx < len(sub); idx++ {
		pos := sub[idx] - 97
		if current.letters[pos] == nil {
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
	for _, pStart := range start.letters {
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
