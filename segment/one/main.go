package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	_ "github.com/mkevac/debugcharts"
	_ "net/http/pprof"

	"github.com/wangriyu/DawnBlossoms/segment/one/version2"
)

const enablePprof = true
var filePath = flag.String("f", "", "files")

func main() {
	flag.Parse()

	if enablePprof {
		go func() {
			// terminal: $ go tool pprof -http=:8081 http://localhost:6060/debug/pprof/heap
			// web:
			// 1、http://localhost:8081/ui
			// 2、http://localhost:6060/debug/charts
			// 3、http://localhost:6060/debug/pprof
			log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		}()
	}

	if len(*filePath) > 0 {
		// testVersion1(*filePath)
		testVersion2(*filePath)
	} else {
		log.Fatalln("file arguments not found")
	}
}

func testVersion2(filepath string) {
	start := time.Now()
	log.Println("v2 start: ", start.String())

	defer func() {
		end := time.Now()
		log.Println("v2 end: ", end.String())
		log.Println("v2 duration: ", end.Sub(start).String())
	}()

	if topHundred, err := version2.CalFinalTopN(filepath); err != nil {
		log.Fatalln(err)
	} else {
		log.Println("v2 result: ", topHundred)
	}
}

// func testVersion1(filepath string) {
// 	start := time.Now()
// 	log.Println("v1 start: ", start.String())
//
// 	defer func() {
// 		end := time.Now()
// 		log.Println("v1 end: ", end.String())
// 		log.Println("v1 duration: ", end.Sub(start).String())
// 	}()
//
// 	if topHundred, err := version1.CalFinalTopN(filepath); err != nil {
// 		log.Fatalln(err)
// 	} else {
// 		log.Println("v1 result: ", topHundred)
// 	}
// }
