/*!
 * <全局变量>
 *
 * Copyright (c) 2019 by <yfm/ ALADING Co.>
 */

package public_yfm

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	MIN_PACKET_LEN   = 4     //数据包的最小长度
	MAX_PACKET_LEN   = 65530 //单个数据包的大小超出这个 就做大包处理
	BIG_PACKET_FLAG  = 65535 //大包标识
	SMALL_PACKET_LEN = 2     //标识小数据包的字节大小(uint16)
	BIG_PACKET_LEN   = 4     //标识大数据包的字节大小(uint32)
	PACKET_CMD_LEN   = 2     //标识数据包命令的字节大小(uint16)
	READ_PACKET_LEN  = 1024  //每次读取的最大字节数
)

const (
	SVR_TYPE_RELATION   = 1 //关系服务器
	SVR_TYPE_CHAT       = 2 //聊天服务器
	SVR_TYPE_GATE       = 3 //网关服务器
	SVR_TYPE_CENTER     = 4 //中心服务器
	SVR_TYPE_GAME       = 5 //游戏服务器
	SVR_TYPE_RECREATION = 6 //娱乐服务器
	SVR_TYPE_TEAMINCOME = 7 //团队收益结算服务器

	DEFAULT_SVRID = 1 //默认服务器id
)

//数据压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()

	return in.Bytes()
}

//全局session计数器
var gSessionID uint32

func GetSessionID() uint32 {
	gSessionID++
	return gSessionID
}

func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	//将[]byte转成16进制
	return fmt.Sprintf("%x", has)
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func CheckPath(path string) {
	_, err := os.Stat(path)
	if err == nil {
		return
	}

	if os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			PrintLog("os.Mkdir error :", err.Error())
		}
		return
	}

	PrintLog("CheckPath error :", err.Error())
}

//一年的第几天
type CISODay struct {
	Year int `bson:"year"`
	Day  int `bson:"day"`
}

//获取外网IP
func GetExternal() (string, error) {
	httpclient := http.Client{Timeout: 2 * time.Second}
	resp, err := httpclient.Get("http://myexternalip.com/raw")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	res, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return "", err2
	}
	r, err3 := regexp.Compile("\\s")
	if err3 != nil {
		return "", err3
	}
	res = r.ReplaceAll(res, []byte(""))
	return string(res), nil
}

//获取空闲端口号
func GetFreePort(num int) ([]int, error) {
	index := 0
	res := make([]int, 0)
	for index < num {
		index++
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			return nil, err
		}
		defer l.Close()
		if v, ok := l.Addr().(*net.TCPAddr); ok {
			res = append(res, v.Port)
		} else {
			return nil, errors.New("*net.tcpaddr assertion failure.")
		}
	}
	return res, nil
}

//获取某天0点时间
// func GetZeroTime(ori time.Time) time.Time {
// 	_, offset := time.Now().Zone()
// 	return time.Unix((ori.Unix()+int64(offset))/SecondsPerDay*SecondsPerDay-int64(offset), 0)
// }

const TimeFormat = "2006-01-02 15:04:05"
const DateFormat = "2006-01-02"

//判断从某个时间点是否跨越了周几
func CheckWeekUpdate(lasttime int64, week int) bool {
	// currweek := int(time.Now().Weekday())
	// if currweek == 0 {
	// 	currweek = 7
	// }

	// if currweek == week {
	// 	return true
	// } else {
	// 	dis := time.Now().Unix() - lasttime
	// 	days := dis / SecondsPerDay
	// 	if days >= 6 {
	// 		return true
	// 	} else {
	// 		tm := time.Unix(lasttime, 0)
	// 		lastweek := int(tm.Weekday())
	// 		if lastweek == 0 {
	// 			lastweek = 7
	// 		}

	// 		if currweek < lastweek { //跨周一的情况
	// 			if currweek > week {
	// 				return true
	// 			}
	// 		} else { //没跨周一的情况
	// 			if currweek > week && lastweek <= week { //上次登录在指定变量week之前 本次登录已过变量week之后
	// 				return true
	// 			}
	// 		}
	// 	}
	// }
	return false
}

//获取程序绝对路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		return strings.Replace(dir, "\\", "/", -1)
	}
	PrintLog("getCurrentDirectory err:", err.Error())
	return ""
}

func GetDayZeroTime() int32 {
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	return int32(tm1.Unix())
}

//获取几天前的0点utc时间
func GetNDayZeroTime(off int) int32 {
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tm2 := tm1.AddDate(0, 0, off)

	return int32(tm2.Unix())
}

//获取随机字符串
//    length：字符串长度
func GetRandomString(length int) string {
	str := "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	var (
		result []byte
		b      []byte
		r      *rand.Rand
	)
	b = []byte(str)
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, b[r.Intn(len(b))])
	}
	return string(result)
}

func GetTodayTimeStr() string {
	currT := time.Now()
	return fmt.Sprintf("%d%d%d", currT.Year(), currT.Month(), currT.Day())
}

func GetYstDayTimeStr() string {
	ystT := time.Now().AddDate(0, 0, -1)
	return fmt.Sprintf("%d%d%d", ystT.Year(), ystT.Month(), ystT.Day())
}
