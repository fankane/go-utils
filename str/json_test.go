package str

import (
	"fmt"
	"testing"
)

var testJ2 = `
[
12.3, false, "hufan"
]
`

var testJ = `
[{
	"key1": 12,
	"key1-1": 12.2,
	"key2": "hello",
	"key3": [
		{
			"key4": true,
			"key5": 3.14
		}
	],
	"hufan": {
		"name": "xx",
		"age": 11,
		"hi": {
			"kk": "hello"
		}
	},
	"score":[1,2,3],
	"score2":[1.3,2.0,3]
}
]
`

func TestParseJSONStr(t *testing.T) {
	res, err := ParseJSONProperty(testJ)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Println(ToJSON(res))
}
