package one

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
)

const (
	splitFileNum = 100 // the split files number
	topN         = 100
)

type splitObj struct {
	file *os.File // file descriptor
	heap *MinHeap // topN minHeap in this file
	list []UrlObj // topN increasing list in this file
}

var (
	regForURL  *regexp.Regexp
	splitFiles [splitFileNum]splitObj
)

func init() {
	var err error
	// regular rules for url
	regForURL, err = regexp.Compile(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	if err != nil {
		panic(err)
	}
	// initialization of split files object
	for i := 0; i < splitFileNum; i++ {
		resultFile, err := os.OpenFile(fmt.Sprintf("split_%d", i), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		splitFiles[i].file = resultFile
		splitFiles[i].heap = NewMinHeap(topN)
		splitFiles[i].list = make([]UrlObj, 0, topN)
	}
}

// splitFile rehash url to one file
// TODO: SplitFile in parallel with goroutines
func SplitFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if matchedAlice := regForURL.FindAll(line, -1); len(matchedAlice) > 0 {
			for _, v := range matchedAlice {
				key := Hash(v)
				if _, err := splitFiles[key].file.Write(append(v, '\n')); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

// calTopNPerFile calculate topN of splitFiles[index].file
// TODO: CalTopNPerFile in parallel with goroutines
func CalTopNPerFile(index int) error {
	var countTable = make(map[string]int, 1<<10) // store url counts。it will take up memory；Escape analysis will assign it to the stack

	splitFiles[index].file.Seek(0, 0) // set seek，otherwise read offset is EOF
	buf := bufio.NewReader(splitFiles[index].file)
	// calculate url counts
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		key := string(line)
		if _, ok := countTable[key]; ok {
			countTable[key]++
		} else {
			countTable[key] = 1
		}
	}

	// push urlObj to minHeap
	for k, v := range countTable {
		if v > splitFiles[index].heap.Top() {
			splitFiles[index].heap.Push(UrlObj{url: k, count: v})
		}
	}

	// transfer minheap to increasing list
	var obj = splitFiles[index].heap.Pop()
	for obj.count > math.MinInt64 {
		// TODO: remove list, just use heap.tree as the list container
		splitFiles[index].list = append(splitFiles[index].list, obj)
		obj = splitFiles[index].heap.Pop()
	}
	// log.Printf("%d %+v %+v\n", index, splitFiles[index].heap, splitFiles[index].list)

	return nil
}
