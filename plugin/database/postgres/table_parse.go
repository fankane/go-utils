package postgres

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/fankane/go-utils/goroutine"
	"github.com/fankane/go-utils/str"
)

type TableColumn struct {
	Idx     int    `json:"idx"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Length  int    `json:"length"`
	Comment string `json:"comment"`
}

const (
	commentDDLPre = "COMMENT ON COLUMN"
)

type ByIdx []*TableColumn

func (a ByIdx) Len() int           { return len(a) }
func (a ByIdx) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIdx) Less(i, j int) bool { return a[i].Idx < a[j].Idx }

func TableColumnsFromDDL(ddl string) ([]*TableColumn, error) {
	lines := strings.Split(ddl, str.LineBreak)
	if len(lines) == 0 {
		return nil, nil
	}
	result := make([]*TableColumn, 0)
	commentLines := make([]string, 0)
	columnFuncs := make([]func() error, 0)
	lock := &sync.Mutex{}
	idx := 1
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !isValidLine(line) {
			continue
		}
		if isComment(line) {
			commentLines = append(commentLines, line)
			continue
		}
		tempLine := line
		tempIdx := idx
		columnFuncs = append(columnFuncs, func() error {
			tc, err := parseSingleColumn(tempLine, tempIdx)
			if err != nil {
				return err
			}
			lock.Lock()
			result = append(result, tc)
			lock.Unlock()
			return nil
		})
		idx++
	}
	if len(columnFuncs) == 0 {
		return nil, fmt.Errorf("empty valided ddl")
	}
	if err := goroutine.Exec(columnFuncs, goroutine.WithReturnWhenError(true)); err != nil {
		return nil, err
	}

	parseColumnsComment(result, commentLines)
	sort.Sort(ByIdx(result))
	return result, nil
}

func isValidLine(line string) bool {
	if line == "" {
		return false
	}
	return isColumnDDL(line) || isComment(line)
}

func isColumnDDL(input string) bool {
	return strings.HasPrefix(input, str.DoubleQuote) && strings.HasSuffix(input, str.Comma)
}

func parseSingleColumn(input string, idx int) (*TableColumn, error) {
	if strings.HasSuffix(input, str.Comma) {
		input = input[:len(input)-1] //去除最后的逗号
	}
	tempArr := strings.Split(input, str.Space)
	if len(tempArr) < 2 {
		return nil, fmt.Errorf("invalid input:[%s]", input)
	}
	tc := &TableColumn{
		Idx:  idx,
		Name: strings.ReplaceAll(tempArr[0], str.DoubleQuote, str.Empty),
		Type: tempArr[1],
	}
	formatColumnLength(tc)
	return tc, nil
}

func isComment(input string) bool {
	return strings.HasPrefix(strings.ToUpper(input), commentDDLPre)
}

// 字段类型里面有 字段长度的提取出来 varchar(255) -> varchar
func formatColumnLength(tc *TableColumn) {
	leftIdx := strings.Index(tc.Type, str.LeftParen)
	if leftIdx < 0 {
		return
	}
	rightIdx := strings.Index(tc.Type, str.RightParen)
	if rightIdx < 0 || rightIdx <= leftIdx {
		return
	}
	middle := tc.Type[leftIdx+1 : rightIdx]
	tc.Type = tc.Type[:leftIdx]
	length, err := strconv.Atoi(middle)
	if err != nil { //非数字 比如 geometry(GEOMETRY)
		return
	}
	tc.Length = length
}

func parseColumnsComment(tcs []*TableColumn, commentLines []string) {
	if len(commentLines) == 0 || len(tcs) == 0 {
		return
	}
	nameComment := make(map[string]string)
	lock := &sync.Mutex{}
	funcList := make([]func() error, 0)
	for _, line := range commentLines {
		tempLine := line
		funcList = append(funcList, func() error {
			name, comment := parseComment(tempLine)
			if name == "" || comment == "" {
				return nil
			}
			lock.Lock()
			nameComment[name] = comment
			lock.Unlock()
			return nil
		})
	}
	if len(funcList) == 0 {
		return
	}
	goroutine.Exec(funcList)
	if len(nameComment) == 0 {
		return
	}
	for _, tc := range tcs {
		if tc == nil {
			continue
		}
		comment, ok := nameComment[tc.Name]
		if ok {
			tc.Comment = comment
		}
	}
}

// COMMENT ON COLUMN "public"."t_result_grid_meta_image_info"."ad_code" IS '行政区划编码'; --> ad_code, 行政区划编码
func parseComment(input string) (string, string) {
	if strings.HasSuffix(input, str.Semicolon) {
		input = input[:len(input)-1] //去除最后的分号
	}
	tempArr := strings.Split(input, str.Space)
	if len(tempArr) < 6 {
		return "", ""
	}
	return getColumnName(tempArr[3]), strings.ReplaceAll(tempArr[5], str.SingleQuote, "")
}

// "public"."t_result_grid_meta_image_info"."ad_code" --> ad_code
func getColumnName(input string) string {
	tempArr := strings.Split(input, str.Dot)
	if len(tempArr) < 3 {
		return ""
	}
	return strings.ReplaceAll(tempArr[2], str.DoubleQuote, str.Empty)
}
