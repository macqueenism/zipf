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

	http.ListenAndServe(":8000", http.DefaultServeMux)
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

	wordMap, count := mapWordsToOccurrence(extractWordsFromLines(readLinesFromFiles(paths)))

	return rankByWordCount(wordMap), count
}

func readLinesFromFiles(paths []string) chan string {
	lineChan := make(chan string, 1)
	done := make(chan bool)
	for _, path := range paths {
		go func(path string) {
			var file *os.File
			file, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			reader := bufio.NewReader(file)

			for {
				str, err := reader.ReadString('\n')
				if err != nil {
					break
				}
				lineChan <- str
			}
			if err == io.EOF {
				err = nil
			}
			done <- true
		}(path)
	}

	go func() {
		for i := 0; i < len(paths); i++ {
			<-done
		}
		close(lineChan)
	}()

	return lineChan
}

func extractWordsFromLines(lineChan chan string) chan string {
	wordChan := make(chan string)
	go func() {
		for line := range lineChan {
			replacer := strings.NewReplacer(",", "", ".", "", ";", "", "?", "", "!", "", "'", "")
			line = replacer.Replace(line)
			splitLine := strings.Split(line, " ")
			for _, word := range splitLine {
				if word != "" && word != "\r\n" {
					wordChan <- word
				}
			}
		}
		close(wordChan)
	}()
	return wordChan
}

func mapWordsToOccurrence(wordChan chan string) (map[string]int, int) {
	wordMap := map[string]int{}
	counter := 0
	for word := range wordChan {
		counter++
		word = strings.ToUpper(word)
		if _, ok := wordMap[word]; ok {
			wordMap[word]++
		} else {
			wordMap[word] = 1
		}
	}
	return wordMap, counter
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
		//fmt.Println(p.Key, ": ", p.Value)
		values = append(values, newVal)
		if i == 20 {
			break
		}
	}
	return values
}
