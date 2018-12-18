package zcmlog

import (
	"log"
	"os"
	"sync"
	"time"
)

var (
	date       string
	file       *os.File
	logger     *log.Logger
	dateFormat = "20060102"
	fileDir    = "./zcmlog/"
	mu         sync.Mutex
)

//第一个参数 文件名字（日期格式） 第二个参数 文件路径
func Init(v ...string) {
	log.Println("【zcmlog init】")
	if len(v) > 0 {
		dateFormat = v[0]
	}
	if len(v) > 1 {
		fileDir = v[1]
	}
	date = time.Now().Format(dateFormat)
	file = createOrOpenFile(date + ".log")
	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
}

func Println(v ...interface{}) {
	if date != time.Now().Format(dateFormat) {
		checkDate()
	}
	//输出到文件，多添加一个换行符
	logger.Println(v, "\n")
}

func checkDate() {
	mu.Lock()
	defer mu.Unlock()
	//如果相等 说明 日期已经处理过
	if date == time.Now().Format(dateFormat) {
		return
	}
	file.Close()
	date = time.Now().Format(dateFormat)
	file = createOrOpenFile(date + ".log")
	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
}
func createOrOpenFile(path string) *os.File {
	os.MkdirAll(fileDir+time.Now().Format("/200601/"), os.ModePerm)
	path = fileDir + time.Now().Format("/200601/") + path
	fi, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	// fi, err := os.Open(path)
	if err != nil {
		if !os.IsExist(err) {
			file, err = os.Create(path)
		}
	}
	if err != nil {
		log.Println("zcmlogerror:" + err.Error())
	}
	return fi
}
