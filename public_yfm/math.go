/*!
 * <字节转换等>
 *
 * Copyright (c) 2019 by <yfm/ ALADING Co.>
 */

package public_yfm

import (
	"bytes"
	"encoding/binary"
	"math"
	"strconv"
	"strings"
)

const (
	YUAN_TO_FEN = 100 //元和分转换系数
)

//4字节数字转字节
func IntToByte(v interface{}) ([]byte, error) {
	header := bytes.NewBuffer([]byte{})
	err := binary.Write(header, binary.LittleEndian, v)
	if err != nil {
		return []byte{}, err
	}
	return header.Bytes(), nil
}

//uint16类型数字
func Uint16(d []byte) uint16 {
	var v uint16 = 0
	header := bytes.NewBuffer(d)
	binary.Read(header, binary.LittleEndian, &v)
	return v
}

//uint32类型数字
func Uint32(d []byte) uint32 {
	var v uint32 = 0
	header := bytes.NewBuffer(d)
	binary.Read(header, binary.LittleEndian, &v)
	return v
}

//int32类型数字
func Int32(d []byte) int32 {
	var v int32 = 0
	header := bytes.NewBuffer(d)
	binary.Read(header, binary.LittleEndian, &v)
	return v
}

//float精度(保留几位小数并四舍五入)
func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

func Round2(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0/pow10_n)*pow10_n) / pow10_n
}

//分转成元
func FenToYuan(value int32) float32 {
	ret := Round(float64(value)/float64(YUAN_TO_FEN), 2)

	return float32(ret)
}

func Int32Slice2Uint16(from []int32) (to []uint16) {
	for _, v := range from {
		to = append(to, uint16(v))
	}
	return
}

func SpliteStrValue(strvar string, spl string) (ret []int32) {
	if strvar == "" {
		return
	}

	value := strings.Split(strvar, spl)
	for _, v1 := range value {
		varl, err := strconv.Atoi(v1)
		if err != nil {
			PrintLog("SpliteStrValue error :", err.Error(), " strvar:", strvar, " spl:", spl)
			panic("SpliteStrValue error")
		}

		ret = append(ret, int32(varl))
	}

	return
}

func SpliteStr2Value(strvar string, sp11, spl2 string) (ret [][]int32) {
	if strvar == "" {
		return
	}

	var1 := strings.Split(strvar, spl2)
	for _, v1 := range var1 {
		value := SpliteStrValue(v1, sp11)
		if len(value) > 0 {
			ret = append(ret, value)
		}
	}

	return
}

const MIN = 0.000001

// MIN 为用户自定义的比较精度
func IsEqual(f1, f2 float64) bool {
	if f1 > f2 {
		return f1-f2 < MIN
	} else {
		return f2-f1 < MIN
	}
}
