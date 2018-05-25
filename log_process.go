package main

import (
	"fmt"
	"strings"
	"time"
	"os"
	"bufio"
	"io"
)

type Reader interface {
	Read(rc chan []byte)
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
func (r *ReadFromFile) Read(rc chan []byte) {
	// 打开文件
	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file err:%s", err.Error()))
	}
	// 移动到文件末尾
	f.Seek(0, 2)

	// 从文件末尾开始逐行读取文件内容
	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500 * time.Microsecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("open file err:%s", err.Error()))
		}
		// 读取模块
		rc <- line[:len(line)-1]
	}
}
func (w WriteToInfluxDB) Write(wc chan string) {
	// 写入模块
	for v := range wc{
		fmt.Println(v)
	}
}

func main() {
	writer := &WriteToInfluxDB{
		dataSource: "username&password",
	}
	reader := &ReadFromFile{
		path: "./access.log",
	}

	lp := &LogProcess{
		rc:     make(chan []byte),
		wc:     make(chan string),
		writer: writer,
		reader: reader,
	}

	go lp.reader.Read(lp.rc)
	go lp.Process()
	go lp.writer.Write(lp.wc)

	time.Sleep(time.Second * 30)
}

type LogProcess struct {
	rc     chan []byte// read chan
	wc     chan string // write chan
	writer Writer
	reader Reader
}

func (l *LogProcess) Process() {
	// 解析模块
	for v := range l.rc {
		l.wc <- strings.ToUpper(string(v))
	}
}
