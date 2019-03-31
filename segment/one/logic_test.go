package one

import (
	"fmt"
	"os"
	"testing"
)

func Test_logic(t *testing.T) {
	tests := []struct {
		name   string
		result []UrlObj
	}{
		{
			name:   "test.log",
			result: []UrlObj{
				{url: "http://wangriyu.wang", count: 27},
				{url: "https://blog.wangriyu.wang", count:26},
				{url: "http://google.com", count:9},
				{url: "https://www.pingcap.com/docs-cn/architecture/", count:1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				for i := 0; i < splitFileNum; i++ {
					if splitFiles[i].file != nil {
						splitFiles[i].file.Close()
					}
					os.Remove(fmt.Sprintf("split_%d", i))
				}
			}()

			// rehash all url string
			err := SplitFile(tt.name)
			if err != nil {
				t.Fatal(err)
			}
			// calculate topN per file
			for i := 0; i < splitFileNum; i++ {
				if err := CalTopNPerFile(i); err != nil {
					t.Fatal(err)
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

			for k, v := range tt.result {
				if topHundred[k].url != v.url || topHundred[k].count != v.count {
					t.Error("unexpected result")
				}
			}
		})
	}
}
