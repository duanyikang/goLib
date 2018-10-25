package main

import (
	"fmt"
	"net/http"
	"net"
	"strings"
	"encoding/json"
	"flag"
	"os"
	"crypto/md5"
	"io"
	"encoding/hex"
	"path/filepath"
)

type PluginFile struct {
	FileName string
	FileMd5  string
	FileUrl  string
}

var fileName string
var fileMd5 string
var fileDownUrl string

func main() {

	filePath := flag.String("f", "", "输入文件的绝对路径")
	flag.Parse()
	fi, err := os.Open(*filePath)
	if err != nil {
		fmt.Println("老哥，你这个文件有问题啊")
		return
	}
	//计算MD5
	md5 := md5.New()
	io.Copy(md5, fi)
	fileMd5 = hex.EncodeToString(md5.Sum(nil))

	fmt.Println("您的文件是:", *filePath)
	fmt.Println("您的文件MD5:", fileMd5)

	paths, filetempName := filepath.Split(*filePath)
	fileName = filetempName
	fileDownUrl = GetPulicIP() + ":8080/" + filetempName
	fmt.Println("您的下载地址:", fileDownUrl)
	StartDownFileServer(filepath.Dir(paths))
	//fileName = "hahah"
	//fileMd5 = "sdfasdfasdfasd"
	//fileDownUrl = "httppp"
	//StartDownFileServer("123")
}

/**
获取本机IP
 */
func GetPulicIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}

/**
开启网络服务
 */
func StartDownFileServer(dir string) {
	http.HandleFunc("/get_plugin_info", get_plugin_info)
	http.Handle("/", http.FileServer(http.Dir(dir)))
	e := http.ListenAndServe(":8080", nil)
	fmt.Println(e)
}

/**
返回插件信息
 */
func get_plugin_info(w http.ResponseWriter, req *http.Request) {
	fmt.Println("loginTask is running...")
	req.ParseForm()

	pluginino := &PluginFile{fileName, fileMd5, "http://"+fileDownUrl}

	//向客户端返回JSON数据
	bytes, _ := json.Marshal(pluginino)
	fmt.Fprint(w, string(bytes))
}
