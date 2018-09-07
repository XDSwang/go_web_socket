// file_socket_Dxx project main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

//连接过的用户数
var i int = 1

//集合 类型[*websocket.Conn]string 用户名不是唯一的
var countryCapitalMap = make(map[*websocket.Conn]string)

//文件路径
var str string = "user_data"

//控制台输出 状态 测试使用 true
var tre bool = true

func Echo(ws *websocket.Conn) {
	//测试使用
	//countryCapitalMap[ws] = "默认名"

	//没发送信息就已经打印 说明连接产生 ws
	WriteWithFileWrite(str, strconv.Itoa(i)+" :用户连接", tre)

	//错误信息
	var err error
	//用户信息
	var reply string

	//更新并发送一次 用户信息
	if err = websocket.Message.Receive(ws, &reply);
	//判断信息不为nil,说明没有返回错误信息.代表成功
	err != nil {
		WriteWithFileWrite(str, countryCapitalMap[ws]+" :用户连接关闭", tre)
		return
	} else {

		//判断是否 设置用户名：
		var flysnowRegexp = regexp.MustCompile(`setname:`)
		var p = flysnowRegexp.FindStringSubmatch(reply)
		//判断p是否匹配到
		if p != nil {
			//是记录或覆盖/更新 集合对信息
			var str_username = strings.Replace(reply, "setname:", "", -1)
			//不允许用户设置空值
			if str_username == "" {
				str_username = "这家伙很懒,设置空名字"
			}
			//加在用户名后面变成唯一,用来给用户识别自己的信息
			//不需要可以在客户端正则去掉
			countryCapitalMap[ws] = str_username + "[" + strconv.Itoa(i) + "]"
			//
			WriteWithFileWrite(str, countryCapitalMap[ws]+" 用户记录成功", tre)

			//发送一次记录的信息给此时的用户
			err = websocket.Message.Send(ws, str_username+"["+strconv.Itoa(i)+"]")
			//发送给用户再修改，不然数据不对
			i++
			//判断是否成功
			if err != nil {

				//证明用户已退出 删除
				delete(countryCapitalMap, ws)

			}

		}

	}

	//连接后用户会一直 循环这个
	for {

		//更新用户信息 			理解：错误会返回 信息，此时err值为nil
		if err = websocket.Message.Receive(ws, &reply);
		//判断信息不为nil,说明没有返回错误信息.代表成功
		err != nil {
			WriteWithFileWrite(str, countryCapitalMap[ws]+" :用户连接关闭", tre)
			break
		}

		//获取存在用户功能 不再解释
		var flysnowRegexp = regexp.MustCompile(`get_users`)
		var p = flysnowRegexp.FindStringSubmatch(reply)
		//判断p是否匹配到
		if p != nil {
			//给集合中存在的用户发送信息 确认用户是否在线的
			for country := range countryCapitalMap {
				//发送信息
				err = websocket.Message.Send(country, "Dxx服务器请求**确认是否在线")

				//判断是否成功
				if err != nil {

					//证明用户已退出 删除
					delete(countryCapitalMap, country)

				}

			}
			err = websocket.Message.Send(ws, "当前在线用户数"+strconv.Itoa(len(countryCapitalMap)))

			//判断是否成功
			if err != nil {

				//证明用户已退出 删除
				delete(countryCapitalMap, ws)

			}
			WriteWithFileWrite(str, countryCapitalMap[ws]+"请求获取用户数,信息："+reply+"当前用户数："+strconv.Itoa(len(countryCapitalMap)), tre)
			continue
		}

		//接收用户发送过来的信息
		WriteWithFileWrite(str, countryCapitalMap[ws]+" 用户发来信息:"+reply, tre)
		//给集合中存在的用户发送信息
		for country := range countryCapitalMap {
			//发送信息
			err = websocket.Message.Send(country, countryCapitalMap[ws]+":"+reply)

			//判断是否成功
			if err != nil {

				//证明用户已退出 删除
				delete(countryCapitalMap, country)

			}

		}

	}
}

//开始函数
func main() {

	//websocket.Handler(Echo)Echo函数 自定义的
	http.Handle("/", websocket.Handler(Echo))
	//监听端口 没有ip默认
	if err := http.ListenAndServe(":1234", nil);
	//失败打印信息
	err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

//功能

//	数据写到文件中 并能控制 控制台是否输出
/*文件路径 写入的字符串 是否控制台输出提示信息 */
func WriteWithFileWrite(name, content string, t bool) {
	fileObj, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {

		控制台输出("打开文件错误:"+err.Error(), t)
		os.Exit(2)
	}
	defer fileObj.Close()
	if _, err := fileObj.WriteString(content + "     时间：" + time.Now().String() + "\n"); err == nil {

		控制台输出("写入成功内容： "+content, t)
	}

}

func 控制台输出(s string, t bool) {
	if t {

		fmt.Println(s)
	}
}
