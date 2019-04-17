package version2

import (
	"testing"
)

func TestLogic(t *testing.T) {
	tests := []struct {
		name   string
		result []UrlObj
	}{
		{
			name:   "test.log",
			result: []UrlObj{
				{Url: "https://blog.wangriyu.wang", Count: 287},
				{Url: "http://blog.wangriyu.wang", Count:246},
				{Url: "https://baidu.com", Count:164},
				{Url: "https://wangriyu.wang", Count:123},
				{Url: "https://google.com", Count:122},
				{Url: "http://wangriyu.wang", Count:121},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if topHundred, err := CalFinalTopN(tt.name); err != nil {
				t.Fatal(err)
			} else {
				t.Log(topHundred)
				for k, v := range tt.result {
					if topHundred[k].Url != v.Url || topHundred[k].Count != v.Count {
						t.Error("unexpected result")
					}
				}
			}
		})
	}
}
