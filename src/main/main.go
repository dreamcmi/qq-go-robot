package main
/*状态码	原因
0	正常
1	错误的auth key
2	指定的Bot不存在
3	Session失效或不存在
4	Session未认证(未激活)
5	发送消息目标不存在(指定对象不存在)
6	指定文件不存在，出现于发送本地图片
10	无操作权限，指Bot没有对应操作的限权
20	Bot被禁言，指Bot当前无法向指定群发送消息
30	消息过长
400	错误的访问，如参数错误等
*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

//基础信息
var url = "http://127.0.0.1:8080"
var qqNum uint32 = 123456

// 创建一个错误处理函数，避免过多的 if err != nil{} 出现
func dropErr(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

//获取版本信息
func getMiraiVersion() string {
	ver, err := http.Get(url + "/about")
	dropErr(err)
	defer ver.Body.Close()
	bs, err := ioutil.ReadAll(ver.Body)
	dropErr(err)

	var miraiVersion gjson.Result = gjson.Get(string(bs), "data.version")
	return miraiVersion.String()
}

//绑定session
func getSessionKey() string{
	//创建发送请求数据
	auth := map[string]string{"authKey": "darrencheng"}
	authJson, err := json.Marshal(auth)
	dropErr(err)

	//开始post
	authGet, err := http.Post(url+"/auth", "application/json", bytes.NewBuffer(authJson))
	defer authGet.Body.Close()
	//读一下返回的body
	authBody, err := ioutil.ReadAll(authGet.Body)
	dropErr(err)
	//gjson解析出session
	var code gjson.Result = gjson.Get(string(authBody), "code")
	if code.Int() == 1 {
		fmt.Println("错误的MIRAI API HTTP auth key")
		return string(authBody)
	}
	var session gjson.Result = gjson.Get(string(authBody), "session")

	/**************************************************************************/
	//创建验证的json
	verify := make(map[string]interface{})
	verify["sessionKey"] = session.String()
	verify["qq"] = qqNum
	verifyJson,err := json.Marshal(verify)
	dropErr(err)

	//开始验证
	verifyGet, err := http.Post(url+"/verify", "application/json", bytes.NewBuffer(verifyJson))
	defer verifyGet.Body.Close()
	//读一下返回的body
	verifyBody, err := ioutil.ReadAll(verifyGet.Body)
	dropErr(err)
	//fmt.Println(string(verifyBody))

	//对返回信息验证
	var verifyReCode gjson.Result = gjson.Get(string(verifyBody),"code")
	if verifyReCode.Int() != 0{
		//fmt.Println(string(verifyBody))
		return string(verifyBody)
	}
	var verifyReMsg gjson.Result = gjson.Get(string(verifyBody),"msg")
	//fmt.Println(verifyReMsg)
	return verifyReMsg.String()
}

func main() {
	ver := getMiraiVersion()
	fmt.Print("mirai_version is:")
	fmt.Println(ver)
	sessionMsg := getSessionKey()
	fmt.Print("getSessionKey:")
	fmt.Println(sessionMsg)
}
