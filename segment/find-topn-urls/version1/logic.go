package version1

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"time"
)

const (
	SplitFileNum = 200 // the split files number
	TopN         = 100
)

type SplitObj struct {
	File *os.File // file descriptor
	Heap *MinHeap // TopN minHeap in this file
	List []UrlObj // TopN increasing list in this file
}

var (
	regForURL  *regexp.Regexp
	SplitFiles [SplitFileNum]SplitObj
)

func init() {
	var err error
	// regular rules for url
	regForURL, err = regexp.Compile(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	if err != nil {
		panic(err)
	}
	// initialization of split files object
	for i := 0; i < SplitFileNum; i++ {
		resultFile, err := os.OpenFile(fmt.Sprintf("split_%d", i), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		SplitFiles[i].File = resultFile
		SplitFiles[i].Heap = NewMinHeap(TopN)
		SplitFiles[i].List = make([]UrlObj, 0, TopN)
	}
}

func CalFinalTopN(filepath string) ([]UrlObj, error) {
	defer func() {
		Close()
	}()

	// rehash all url string
	err := SplitFile(filepath)
	if err != nil {
		return nil, err
	}

	t1 := time.Now()
	// calculate topN per file
	for i := 0; i < SplitFileNum; i++ {
		if err := CalTopNPerFile(i); err != nil {
			return nil, err
		}
	}

	// calculate final topN
	var topN = make([]UrlObj, 0, TopN)
	for len(topN) < TopN {
		max := UrlObj{}
		index := 0
		for i := 0; i < SplitFileNum; i++ {
			len := len(SplitFiles[i].List)
			if len > 0 && SplitFiles[i].List[len-1].Count > max.Count {
				max = SplitFiles[i].List[len-1]
				index = i
			}
		}
		if l := len(SplitFiles[index].List); l > 0 {
			SplitFiles[index].List = SplitFiles[index].List[:l-1]
		}
		topN = append(topN, max)
	}
	log.Println("v1 calTopN: ", time.Now().Sub(t1).String())

	return topN, nil
}

// splitFile rehash url to find-topn-urls file
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
				key := Hash(v) % SplitFileNum
				if _, err := SplitFiles[key].File.Write(append(v, '\n')); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

// calTopNPerFile calculate TopN of SplitFiles[index].file
// TODO: CalTopNPerFile in parallel with goroutines
func CalTopNPerFile(index int) error {
	var countTable = make(map[string]int, 1<<10) // store url counts。it will take up memory；Escape analysis will assign it to the stack

	SplitFiles[index].File.Seek(0, 0) // set seek，otherwise read offset is EOF
	buf := bufio.NewReader(SplitFiles[index].File)
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
		if SplitFiles[index].Heap.len < SplitFiles[index].Heap.max || v > SplitFiles[index].Heap.Top() {
			SplitFiles[index].Heap.Push(UrlObj{Url: k, Count: v})
		}
	}

	// transfer minheap to increasing list
	var obj = SplitFiles[index].Heap.Pop()
	for obj.Count > math.MinInt64 {
		// TODO: remove list, just use heap.tree as the list container
		SplitFiles[index].List = append(SplitFiles[index].List, obj)
		obj = SplitFiles[index].Heap.Pop()
	}

	return nil
}

func Close() {
	for i := 0; i < SplitFileNum; i++ {
		if SplitFiles[i].File != nil {
			SplitFiles[i].File.Close()
		}
		os.Remove(fmt.Sprintf("split_%d", i))
	}
}
