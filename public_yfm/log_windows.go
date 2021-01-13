/*!
 * <日志记录>
 *
 * Copyright (c) 2019 by <yfm/ ALADING Co.>
 */

package public_yfm

import (
	"fmt"
	"log"
	"os"
	"time"
)

var gFile *os.File

//记录日志到文件(每日一个文件)并打印到控制台
func PrintLog(v ...interface{}) {
	fileName := "./log/" + fmt.Sprintf("%d_%d_%d.log", time.Now().Year(), time.Now().Month(), time.Now().Day())
	if gFile == nil || gFile.Name() != fileName {
		if gFile != nil {
			gFile.Close()
		}
		var err error
		gFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
		if err != nil {
			log.Println("Log Open File Error :", err.Error())
			return
		}
	}
	pLoger := log.New(gFile, "", log.LstdFlags)
	log.Println(fmt.Sprint(v...))
	pLoger.Println(fmt.Sprint(v...))
}

//记录日志到指定的文件并打印到控制台
func PrintFileLog(filename string, v ...interface{}) {
	file, err := os.OpenFile("./log/"+filename+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		log.Println("FileLog Error :", filename, " Error:", err.Error())
		return
	}

	pLoger := log.New(file, "", log.LstdFlags)
	log.Println(fmt.Sprint(v...))
	pLoger.Println(fmt.Sprint(v...))
	file.Close()
}
