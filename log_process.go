package main

import (
	"fmt"
	"strings"
	"time"
)

type Reader interface {
	Read(rc chan string)
}

type Writer interface {
	Write(wc chan string)
}

type ReadFromFile struct {
	path string
}
type WriteToInfluxDB struct {
	dataSource string
}

// 使用引用的话可以不用发生值拷贝
// 并且如果需要修改这个结构体本身的一些属性时，可以直接使用l变量去修改
func (r *ReadFromFile) Read(rc chan string) {
	// 读取模块
	rc <- "message"
}
func (w WriteToInfluxDB) Write(wc chan string) {
	// 写入模块
	fmt.Println(<-wc)
}

func main() {
	writer := &WriteToInfluxDB{
		dataSource: "username&password",
	}
	reader := &ReadFromFile{path: "path"}


	lp := &LogProcess{
		rc:   make(chan string),
		wc:   make(chan string),
		writer:writer,
		reader:reader,
	}

	go lp.reader.Read(lp.rc)
	go lp.Process()
	go lp.writer.Write(lp.wc)

	time.Sleep(time.Second * 1)
}

type LogProcess struct {
	rc     chan string // read chan
	wc     chan string // write chan
	writer Writer
	reader Reader
}

func (l *LogProcess) Process() {
	// 解析模块j
	data := <-l.rc
	l.wc <- strings.ToUpper(data)
}
