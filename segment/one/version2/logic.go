package version2

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"runtime"
	"sync"
	"time"
)

const (
	SplitFileNum = 20 // the split files number
	TopN         = 10
)

type SplitObj struct {
	File    *os.File // file descriptor
	Heap    *MinHeap // TopN minHeap in this file
	WriteCh chan string
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
		SplitFiles[i].WriteCh = make(chan string, 1<<12)
	}
}

func CalFinalTopN(filepath string) ([]UrlObj, error) {
	defer func() {
		Close()
		// if err := recover(); err != nil {
		// 	log.Println(err)
		// }
	}()

	// rehash all url string
	SplitFile(filepath)

	t1 := time.Now()
	errCh := make(chan error)
	doneCh := make(chan bool)

	go parallelCalTopN(doneCh, errCh)

	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
		return nil, nil
	case <-doneCh:
		// calculate final topN
		var topN = make([]UrlObj, 0, TopN)
		for len(topN) < TopN {
			max := UrlObj{}
			index := 0
			for i := 0; i < SplitFileNum; i++ {
				l := len(SplitFiles[i].Heap.tree)
				if l > 0 && SplitFiles[i].Heap.tree[0].Count > max.Count {
					max = SplitFiles[i].Heap.tree[0]
					index = i
				}
			}
			if l := len(SplitFiles[index].Heap.tree); l > 0 {
				SplitFiles[index].Heap.tree = SplitFiles[index].Heap.tree[1:]
			}
			topN = append(topN, max)
		}
		log.Println("v2 calTopN: ", time.Now().Sub(t1).String())

		return topN, nil
	}
}

// SplitFile rehash url to one file
func SplitFile(filepath string) {
	eofCh := make(chan bool)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go produceMsg(cancel, filepath, eofCh)

	wg := sync.WaitGroup{}
	for i := 0; i < SplitFileNum; i++ {
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, index int, eofCh chan bool) {
			defer func() {
				wg.Done()
				// if err := recover(); err != nil {
				// 	log.Println(err)
				// 	if e , ok := err.(error); ok {
				// 		errCh <- e
				// 	}
				// }
			}()
			for {
				select {
				case <-eofCh:
					if len(SplitFiles[index].WriteCh) == 0 {
						return
					}
				case <-ctx.Done():
					return
				case data, ok := <-SplitFiles[index].WriteCh:
					if !ok {
						log.Printf("split file %d writeCh closed\n", index)
						return
					}
					if _, err := SplitFiles[index].File.WriteString(data); err != nil {
						panic(err)
						return
					}
				}
			}
		}(ctx, &wg, i, eofCh)
	}
	wg.Wait()
}

func produceMsg(cancel context.CancelFunc, filepath string, eofCh chan bool) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		log.Println(err)
	// 		if e , ok := err.(error); ok {
	// 			errCh <- e
	// 		}
	// 		cancel()
	// 	}
	// }()

	file, err := os.Open(filepath)
	if err != nil {
		log.Println(err)
		panic(err)
		return
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			log.Println(err)
			if err == io.EOF {
				close(eofCh)
				return
			}
			panic(err)
			cancel()
			return
		}
		if matchedAlice := regForURL.FindAll(line, -1); len(matchedAlice) > 0 {
			for _, v := range matchedAlice {
				key := Hash(v) % SplitFileNum
				SplitFiles[key].WriteCh <- fmt.Sprintf("%s\n", v)
			}
		}
	}
}

// calculate topN per file in parallel
func parallelCalTopN(doneCh chan bool, errCh chan error) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		if e, ok := err.(error); ok {
	// 			errCh <- e
	// 		}
	// 	}
	// }()
	for i := 0; i < SplitFileNum; {
		num := runtime.NumCPU()
		if num > 4 {
			num = 4
		}
		wg := sync.WaitGroup{}
		for j := 0; j < num && i+j < SplitFileNum; j++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, index int) {
				defer wg.Done()
				if err := CalTopNPerFile(index); err != nil {
					errCh <- err
				}
			}(&wg, i+j)
		}
		wg.Wait()
		i += num
	}
	close(doneCh)
}

// CalTopNPerFile calculate TopN of SplitFiles[index].file
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

	// transfer minheap.tree to increasing list
	var obj = SplitFiles[index].Heap.PopToList()
	for obj.Count > math.MinInt64 {
		obj = SplitFiles[index].Heap.PopToList()
	}
	SplitFiles[index].Heap.tree = SplitFiles[index].Heap.tree[1:]
	// log.Printf("calTopNPeFile %d %+v\n", index, SplitFiles[index].Heap)

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
