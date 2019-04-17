package main

import (
	"bytes"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

const characters = "abcdefghijklmnopqrstuvwxyz.0123456789"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	file, err := os.OpenFile("test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	max := (1 << 28)/runtime.NumCPU()
	wg := sync.WaitGroup{}
	for i := 0; i <= runtime.NumCPU(); i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, num int, file *os.File) {
			defer func() {
				wg.Done()
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()
			log.Printf("goroutine %d start...\n", num)
			for j := 0; j <= max; j++ {
				writeUrl(file)
			}
		}(&wg, i, file)
	}
	wg.Wait()
}

func writeUrl(file *os.File) {
	var bf []byte
	str := bytes.NewBuffer(bf)
	str.Grow(25)
	str.WriteString("http://www.")
	str.Write(RandString(10))
	str.WriteString(".com\n")
	if _, err := file.Write(str.Bytes()); err != nil {
		log.Println(err)
	}
}

func RandString(n int) []byte {
	b := make([]byte, n)
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := int64(len(characters))
	for i := range b {
		b[i] = characters[rd.Int63n(l)]
	}
	return b
}
