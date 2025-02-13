package snowflake

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/fankane/go-utils/str"
	"github.com/fankane/go-utils/utime"
)

func TestNewSnowflake(t *testing.T) {
	sg := sync.WaitGroup{}
	sg.Add(1)
	go func() {
		defer sg.Done()
		s, _ := NewSnowflake(1, 1)

		for i := 0; i < 3; i++ {
			id, err := s.NextID()
			if err != nil {
				fmt.Println("err:", err)
				return
			}
			fmt.Println("s1:", id)
			time.Sleep(time.Millisecond * 2)
		}
	}()
	sg.Add(1)
	go func() {
		defer sg.Done()
		st, _ := time.Parse(utime.LayYMDHms1, "2024-02-13 09:49:23")
		s2, _ := NewSnowflake(1, 1, CustomizeEpoch(st))
		for i := 0; i < 3; i++ {
			id, err := s2.NextID()
			if err != nil {
				fmt.Println("err:", err)
				return
			}
			fmt.Println("s2:", id)
			time.Sleep(time.Millisecond * 2)
		}
	}()
	sg.Wait()
}

func TestSnowflake_ParseSnowflakeID(t *testing.T) {
	s, _ := NewSnowflake(1, 12)

	for i := 0; i < 3; i++ {
		id, err := s.NextID()
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		fmt.Println("s1:", id)
		pInfo, err := s.ParseSnowflakeID(id)
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		fmt.Println(str.ToJSON(pInfo))
		//time.Sleep(time.Second * 2)
	}
}
