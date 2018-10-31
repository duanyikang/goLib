```
package main

import (
	"fmt"
	"net"
	"strings"
	"flag"
	"os"
	"crypto/md5"
	"io"
	"encoding/hex"
	"path/filepath"
	"net/http"
	"encoding/json"
)

type PluginFile struct {
	FileName string
	FileMd5  string
	FileUrl  string
}

var filePath_gt string //各推的路径
var filePath_hw string //华为的路径
var filePath_xm string //小米的路径

var fileMd5_hw string
var fileMd5_xm string
var fileMd5_gt string

var fileDownUrl_hw string
var fileDownUrl_xm string
var fileDownUrl_gt string

var pluginfo PluginFile

func main() {

	filePath_hw := flag.String("hw", "/Users/duanyikang/Desktop/com.yichuizi.pushplugin_20181030181038.apk", "华为输入文件的绝对路径")
	fileName_xm := flag.String("xm", "/Users/duanyikang/Desktop/河豚钥匙.apk", "小米输入文件的绝对路径")
	fileName_gt := flag.String("gt", "/Users/duanyikang/Desktop/channel.txt", "个推输入文件的绝对路径")
	flag.Parse()
	checkPlugin(*filePath_hw, 1)
	checkPlugin(*fileName_xm, 2)
	checkPlugin(*fileName_gt, 3)

	StartDownFileServer()
}

func checkPlugin(path string, t int) {
	fi, err := os.Open(path)
	if err != nil {
		fmt.Println("老哥，你这个文件有问题啊")
		return
	}
	//计算MD5
	md5 := md5.New()
	io.Copy(md5, fi)

	paths, filetempName := filepath.Split(path)

	switch t {
	case 1:
		fileMd5_hw = hex.EncodeToString(md5.Sum(nil))
		fileDownUrl_hw = GetPulicIP() + ":8080/" + filetempName
		filePath_hw = paths
		break
	case 2:
		fileMd5_xm = hex.EncodeToString(md5.Sum(nil))
		fileDownUrl_xm = GetPulicIP() + ":8080/" + filetempName
		filePath_xm = paths
		break
	case 3:
		fileMd5_gt = hex.EncodeToString(md5.Sum(nil))
		fileDownUrl_gt = GetPulicIP() + ":8080/" + filetempName
		filePath_gt = paths
		break
	}

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
func StartDownFileServer() {
	http.HandleFunc("/get_plugin_info", get_plugin_info)
	fmt.Println("我要的："+filePath_hw)
	http.Handle("/", http.FileServer(http.Dir(filePath_hw)))
	//http.Handle("/", http.FileServer(http.Dir(filePath_xm)))
	//http.Handle("/", http.FileServer(http.Dir(filePath_gt)))
	e := http.ListenAndServe(":8080", nil)
	fmt.Println(e)
}

/**
返回插件信息
 */
func get_plugin_info(w http.ResponseWriter, req *http.Request) {
	fmt.Println("有客户端访问我了")
	deviceName := req.FormValue("deviceName")

	fmt.Println("我要的参数:", deviceName)

	if strings.Contains(deviceName, "xiaomi") {
		fmt.Println("小米的设备")
		pluginfo = PluginFile{fileDownUrl_xm, fileMd5_xm, "http://" + fileDownUrl_xm}
	} else if strings.Contains(deviceName, "huawei") {
		fmt.Println("华为的设备")
		pluginfo = PluginFile{fileDownUrl_hw, fileMd5_hw, "http://" + fileDownUrl_hw}
	} else {
		fmt.Println("其他的的设备")
		pluginfo = PluginFile{fileDownUrl_gt, fileMd5_gt, "http://" + fileDownUrl_gt}
	}
	//向客户端返回JSON数据
	bytes, _ := json.Marshal(pluginfo)
	fmt.Println("您的下载地址:",  pluginfo.FileUrl)
	fmt.Println("您的MD5:",pluginfo.FileMd5)
	fmt.Fprint(w, string(bytes))
}


```
