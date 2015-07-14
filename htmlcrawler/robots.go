package htmlcrawler

import (
	"bufio"
	"errors"
	"net/http"
	"strings"
)

//TODO: Use path to join paths instead of concat
func LoadRobots(root string) (map[string]bool, []string, error) {
	disallows := make(map[string]bool)
	allows := make([]string, 0)
	resp, err := http.Get(root + "/robots.txt")
	if err != nil {
		return disallows, allows, errors.New("Error loading robots")
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		information := strings.SplitN(scanner.Text(), ":", 2)
		if len(information) == 2 {
			switch information[0] {
			case "Sitemap":
				allows = append(allows, strings.TrimSpace(information[1]))
			case "Disallow":
				disallowedUrl := root + strings.TrimSpace(information[1])
				disallows[disallowedUrl] = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return disallows, allows, errors.New("Error reaching end of robots")
	}
	return disallows, allows, nil
}
