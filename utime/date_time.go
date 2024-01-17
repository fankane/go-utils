package utime

import (
	"time"
)

const (
	LayY       = "2006" //年
	LayM       = "01"   //月
	LayD       = "02"   //日
	LayH       = "15"   //时
	LayMin     = "04"   //分
	Lays       = "05"   //秒
	LayYMD1    = "2006-01-02"
	LayYMD2    = "2006/01/02"
	LayYMDHms1 = "2006-01-02 15:04:05"
	LayYMDHms2 = "2006/01/02 15:04:05"
	LayYMDHms3 = "20060102150405"
	LayHms     = "15:04:05"
)

func GetUTC8Time() time.Time {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}
	}
	return time.Now().In(location)
}

func GetUTC8Loc() *time.Location {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil
	}
	return location
}
