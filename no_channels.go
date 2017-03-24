package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/wcharczuk/go-chart"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {

	start := time.Now()

	pairList, total := countWordDist(getTxtFiles("./books/"))

	elapsed := time.Since(start)
	fmt.Println("total time elapsed: ", elapsed)
	fmt.Println("total number of words: ", total)

	generateChart(pairList, total)

	http.ListenAndServe(":8080", http.DefaultServeMux)
}

func getTxtFiles(path string) (fileList []string) {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if strings.Contains(f.Name(), ".txt") {
			fileAndDir := path + f.Name()
			fileList = append(fileList, fileAndDir)
		}
	}
	return
}

func countWordDist(paths []string) (PairList, int) {
	totalLines := []string{}
	for _, path := range paths {
		lines, err := readLines(path)
		if err != nil {
			fmt.Println("readLines: ", err)
			continue
		}
		for _, line := range lines {
			totalLines = append(totalLines, line)
		}
	}

	wordList := extractWordsFromLines(totalLines)

	wordMap := mapWordsToOccurrence(wordList)

	return rankByWordCount(wordMap), len(wordList)
}

func readLines(path string) (lines []string, err error) {
	var (
		file *os.File
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		lines = append(lines, str)
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func extractWordsFromLines(totalLines []string) []string {
	var wordList []string
	for _, line := range totalLines {
		replacer := strings.NewReplacer(",", "", ".", "", ";", "", "?", "", "!", "", "'", "")
		line = replacer.Replace(line)
		splitLine := strings.Split(line, " ")
		for _, word := range splitLine {
			if word != "" && word != "\r\n" {
				wordList = append(wordList, word)
			}
		}
	}
	return wordList
}

func mapWordsToOccurrence(wordList []string) map[string]int {
	wordMap := map[string]int{}
	for _, word := range wordList {
		word = strings.ToUpper(word)
		if _, ok := wordMap[word]; ok {
			wordMap[word]++
		} else {
			wordMap[word] = 1
		}
	}
	return wordMap
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func generateChart(pairList PairList, total int) {
	graph := chart.BarChart{
		Height:   2000,
		BarWidth: 100,
		XAxis: chart.Style{
			Show: true,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Bars: buildBarValues(pairList, total),
	}

	file, _ := os.Create("chart.png")
	defer file.Close()

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, buffer)

	file.Write(buffer.Bytes())
}

func buildBarValues(pairList PairList, total int) []chart.Value {
	values := []chart.Value{}
	for i, p := range pairList {
		newVal := chart.Value{
			Value: float64(p.Value) / float64(total) * 100,
			Label: p.Key,
		}
		values = append(values, newVal)
		if i == 20 {
			break
		}
	}
	return values
}
