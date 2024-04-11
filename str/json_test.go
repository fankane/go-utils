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
  {
    "Center": [
      116.410503,
      39.911502
    ],
    "Children": [
	{
		"id": "101"
	}],
    "Code": "110101",
    "Level": "district",
    "Name": "东城区",
    "bound": {}
  }
`

func TestParseJSONStr(t *testing.T) {
	res, err := ParseJSONProperty(testJ)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Println(ToJSON(res))
}
