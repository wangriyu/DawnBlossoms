package main

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
		splitFiles[i].heap = NewMinheap(topN)
		splitFiles[i].list = make([]UrlObj, 0, topN)
	}
}

func main() {
	defer func() {
		for i := 0; i < splitFileNum; i++ {
			splitFiles[i].file.Close()
			os.Remove(fmt.Sprintf("split_%d", i))
		}
	}()

	// rehash all url string
	err := splitFile("test.log")
	if err != nil {
		log.Panic(err)
	}
	// calculate topN per file
	for i := 0; i < splitFileNum; i++ {
		if err := calTopNPerFile(i); err != nil {
			log.Panic(err)
		}
	}

	// calculate final topN
	var topHundred = make([]UrlObj, 0, topN)
	for len(topHundred) < topN {
		max := UrlObj{}
		index := 0
		for i := 0; i < splitFileNum; i++ {
			len := len(splitFiles[i].list)
			if len > 0 && splitFiles[i].list[len-1].count > max.count {
				max = splitFiles[i].list[len-1]
				index = i
			}
		}
		if l := len(splitFiles[index].list); l > 0 {
			splitFiles[index].list = splitFiles[index].list[:l-1]
		}
		topHundred = append(topHundred, max)
	}
	fmt.Printf("%+v", topHundred)
}

// splitFile rehash url to one file
func splitFile(filepath string) error {
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
				key := hash(v)
				if _, err := splitFiles[key].file.Write(append(v, '\n')); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

// calTopNPerFile calculate topN of splitFiles[index].file
func calTopNPerFile(index int) error {
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
		splitFiles[index].list = append(splitFiles[index].list, obj)
		obj = splitFiles[index].heap.Pop()
	}
	// log.Printf("%d %+v %+v\n", index, splitFiles[index].heap, splitFiles[index].list)

	return nil
}

type UrlObj struct {
	url   string
	count int
}

type MinHeap struct {
	max  int
	len  int
	tree []UrlObj
}

func NewMinheap(max int) *MinHeap {
	heap := &MinHeap{max: max, tree: make([]UrlObj, 1, max)}
	heap.tree[0] = UrlObj{count: math.MinInt64}
	return heap
}

func (heap *MinHeap) Push(x UrlObj) {
	if heap.len >= heap.max {
		heap.Pop()
	}
	heap.tree = append(heap.tree, x)
	heap.len++
	i := heap.len
	for ; heap.tree[i/2].count > x.count; i /= 2 {
		heap.tree[i] = heap.tree[i/2]
	}
	heap.tree[i] = x
}

func (heap *MinHeap) Top() int {
	if heap.len <= 0 {
		return heap.tree[0].count
	}
	return heap.tree[1].count
}

func (heap *MinHeap) Pop() UrlObj {
	if heap.len > 0 {
		min := heap.tree[1]
		last := heap.tree[heap.len]
		var i, child int
		for i = 1; i*2 <= heap.len; i = child {
			child = i * 2
			if child < heap.len && heap.tree[child+1].count < heap.tree[child].count {
				child++
			}
			if last.count > heap.tree[child].count {
				heap.tree[i] = heap.tree[child]
			} else {
				break
			}
		}
		heap.tree[i] = last
		heap.tree = heap.tree[:heap.len]
		heap.len--
		return min
	}
	return heap.tree[0]
}

func hash(str []byte) (key int) {
	for _, v := range str {
		key += int(v)
	}
	return key % splitFileNum
}
