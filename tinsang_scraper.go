package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"sync"
)

const NumTinSangArticles = 3101

type mapSorter struct {
	original_map map[string]int
	keys         []string
}

func (ms *mapSorter) Len() int {
	return len(ms.original_map)
}

func (ms *mapSorter) Less(i, j int) bool {
	return ms.original_map[ms.keys[i]] > ms.original_map[ms.keys[j]]
}

func (ms *mapSorter) Swap(i, j int) {
	ms.keys[i], ms.keys[j] = ms.keys[j], ms.keys[i]
}

func GetSortedMapKeysByValue(m map[string]int) []string {
	ms := new(mapSorter)
	ms.original_map = m
	ms.keys = make([]string, len(m))
	i := 0
	for key, _ := range m {
		ms.keys[i] = key
		i++
	}
	sort.Sort(ms)
	return ms.keys
}

func NewStringSet() *StringSet {
	return &StringSet{make(map[string]bool)}
}

type StringSet struct {
	items map[string]bool
}

func (set *StringSet) Add(s string) {
	set.items[s] = true
}

func (set *StringSet) Remove(s string) {
	delete(set.items, s)
}

func ParseUserScoreFromUserName(user_name string) (int, error) {
	response, err := http.Get("http://tinsang.net/user/" + user_name)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}
	re := regexp.MustCompile("<b>quả đã có </b>([0-9]+)")
	matches := re.FindAllStringSubmatch(string(contents), -1)
	if len(matches) == 0 {
		return 0, errors.New("Cannot find score for user: " + user_name)
	}
	score, _ := strconv.Atoi(matches[0][1])
	return score, nil
}

// Get user strings of form /user/foo.
func ParseUserStringsFromTinSangUrl(url string) []string {
	response, err := http.Get(url)
	if err != nil {
		return make([]string, 0)
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return make([]string, 0)
	}
	re := regexp.MustCompile("/user/[^\"]*")
	return re.FindAllString(string(contents), -1)
}

func GenUsersNamesFromTinSangUrl(url string) <-chan string {
	out := make(chan string)
	go func() {
		for _, user_string := range ParseUserStringsFromTinSangUrl(url) {
			user_name := user_string[6:len(user_string)]
			out <- user_name
		}
		close(out)
	}()
	return out
}

func GetAllTinSangUserNames(user_names *StringSet) {
	var wg sync.WaitGroup
	out := make(chan string)
	output := func(c <-chan string) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(NumTinSangArticles)
	for article_index := 1; article_index <= NumTinSangArticles; article_index++ {
		article_url := "http://tinsang.net/news/" +
			strconv.Itoa(article_index)
		go output(GenUsersNamesFromTinSangUrl(article_url))
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	for name := range out {
		user_names.Add(name)
	}
}

func main() {
	user_names := NewStringSet()
	user_names_to_points := make(map[string]int)
	GetAllTinSangUserNames(user_names)
	for name := range user_names.items {
		score, err := ParseUserScoreFromUserName(name)
		if err == nil {
			user_names_to_points[name] = score
		}
	}
	sorted_names_by_scores := GetSortedMapKeysByValue(user_names_to_points)
	for _, name := range sorted_names_by_scores {
		fmt.Println(name + ":" + strconv.Itoa(user_names_to_points[name]))
	}
}
