// 河北金融学院 OA 疫情数据自动填报（新版-适配微信小程序版本）
// Powered By Luckykeeper <luckykeeper@luckykeeper.site | https://luckykeeper.site> 2022/09/06
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	strftime "github.com/jehiah/go-strftime"
	"github.com/tidwall/gjson"
)

var (
	oaUsername string // OA 账号
	oaPassword string // OA 密码
	address    string // 详细地址
	prove      bool   // 核酸检测证明，有填 true ，没有填 false
	addressGA  string // GitHub Action 用地址，URL编码

	addressGAToStr, HaveProve, AuthorizationCodeStr string
)

func init() {
	flag.StringVar(&oaUsername, "oaUsername", "", "河北金融学院OA系统的用户名(学工号)")
	flag.StringVar(&oaPassword, "oaPassword", "", "河北金融学院OA系统的密码")
	flag.StringVar(&address, "address", "", "当前的居住地址(省/市/区/街道/详细地址),省市参数可以在 https://oa.hbfu.edu.cn/datafill/collect/usertask 获取")
	flag.BoolVar(&prove, "prove", false, "是否持有核酸证明,可选(true|false),true=>有,false=>无")
	flag.StringVar(&addressGA, "addressGA", "", "当前的居住地址Url编码（仅供GitHub Action使用）")
	flag.Parse()
	addressGAToStr, _ = url.QueryUnescape(addressGA)
}

func main() {

	fmt.Println("河北金融学院自动每日健康打卡（新版）Ver1.01 Build20220906")
	fmt.Println("Powered By Luckykeeper <luckykeeper@luckykeeper.site | https://luckykeeper.site>")
	fmt.Println("GitHub:https://github.com/HBFUer/dataFill_new")
	fmt.Println("_____________________________________________")
	if oaUsername == "" || oaPassword == "" { // 校验：OA用户名密码必须填写
		usage()
	} else if addressGA == "" && address == "" { // 校验，两个地址不能都空
		usage()
	} else if addressGA != "" && address != "" { // 校验，两个地址不能都填
		usage()
	} else {
		if addressGA != "" {
			address = addressGAToStr
		}
		fmt.Println("本程序将自动完成每日健康打卡，你需要对你上报的数据负责！程序仅负责调用接口上报数据！")
		fmt.Println("程序仅供学习探讨Go语言编程，对使用本程序造成的一切后果作者均不负责！")
		fmt.Println("程序不存储用户账户密码，请妥善保管好相关信息！")
		fmt.Println("程序不对接口变动后可能产生的异常负责，请关于接口信息！")
		fmt.Println("运行程序则代表已知晓并同意以上规则！")
		createDir("./Screenshots")
		getAuthAndDataFill()
	}
}

// Usage
func usage() {
	fmt.Println("Usage:")
	fmt.Println("-oaUsername (必选参数)河北金融学院OA系统的用户名(学工号)")
	fmt.Println("-oaPassword (必选参数)河北金融学院OA系统的密码")
	fmt.Println("-address (必选参数，和 addressGA 必选其中一个)当前的居住地址(省/市/区/街道/详细地址),省市参数可以在 https://oa.hbfu.edu.cn/datafill/collect/usertask 获取")
	fmt.Println("-addressGA (必选参数，和 address 必选其中一个)当前的居住地址Url编码（仅供GitHub Action使用）")
	fmt.Println("-prove (可选参数)是否持有核酸证明,可选(true|false),true=>有,false=>无。默认：无")
}

// 获取 AuthorizationCode
func getAuthAndDataFill() {

	var pic1 []byte // debug 使用
	var pic0, pic2 []byte
	// create context
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true)) // debug(false)|prod(true)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(
		allocCtx,
	)
	defer cancel()

	getCodeAndDataFill(ctx)

	if err := chromedp.Run(ctx,
		// OA酱 登录页配置
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate("https://oa.hbfu.edu.cn/backstage/cas/login?service=https%3A%2F%2Foa.hbfu.edu.cn%2Fbackstage%2Fcas-proxy%2Fapp%2Fredirect"),
		chromedp.WaitVisible(`body > app-root > app-theme1 > div > main > div.right-container.box-shadow > app-login-panel > div > nz-tabset > div.ant-tabs-content.ant-tabs-top-content.ant-tabs-content-animated > div.ant-tabs-tabpane.ant-tabs-tabpane-active > app-userpass-panel > form > input:nth-child(1)`),
		chromedp.Sleep(2*time.Second),
		chromedp.FullScreenshot(&pic0, 90),
		chromedp.Click(`body > app-root > app-theme1 > div > main > div.right-container.box-shadow > app-login-panel > div > nz-tabset > div.ant-tabs-content.ant-tabs-top-content.ant-tabs-content-animated > div.ant-tabs-tabpane.ant-tabs-tabpane-active > app-userpass-panel > form > input:nth-child(1)`, chromedp.ByQuery),
		chromedp.SendKeys(`body > app-root > app-theme1 > div > main > div.right-container.box-shadow > app-login-panel > div > nz-tabset > div.ant-tabs-content.ant-tabs-top-content.ant-tabs-content-animated > div.ant-tabs-tabpane.ant-tabs-tabpane-active > app-userpass-panel > form > input:nth-child(1)`, oaUsername, chromedp.ByQuery),
		chromedp.Click(`body > app-root > app-theme1 > div > main > div.right-container.box-shadow > app-login-panel > div > nz-tabset > div.ant-tabs-content.ant-tabs-top-content.ant-tabs-content-animated > div.ant-tabs-tabpane.ant-tabs-tabpane-active > app-userpass-panel > form > input:nth-child(2)`, chromedp.ByQuery),
		chromedp.SendKeys(`body > app-root > app-theme1 > div > main > div.right-container.box-shadow > app-login-panel > div > nz-tabset > div.ant-tabs-content.ant-tabs-top-content.ant-tabs-content-animated > div.ant-tabs-tabpane.ant-tabs-tabpane-active > app-userpass-panel > form > input:nth-child(2)`, oaPassword, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.SendKeys(`body > app-root > app-theme1 > div > main > div.right-container.box-shadow > app-login-panel > div > nz-tabset > div.ant-tabs-content.ant-tabs-top-content.ant-tabs-content-animated > div.ant-tabs-tabpane.ant-tabs-tabpane-active > app-userpass-panel > form > input:nth-child(2)`, kb.Enter, chromedp.ByQuery),
		chromedp.FullScreenshot(&pic1, 90),
		chromedp.Sleep(2*time.Second),
		// OA酱 - 跳转到首页完成登录流程
		chromedp.Navigate("https://oa.hbfu.edu.cn/new/angular-office-hall/#/angular-office-hall/index"),
		chromedp.WaitVisible(`#hallBody > app-root > app-index > nz-layout > app-header > div > div > div:nth-child(2) > div > div > div > ul > li:nth-child(2)`, chromedp.ByQuery),
		// OA酱 - 进入数据填报
		chromedp.Navigate("https://oa.hbfu.edu.cn/backstage/mars-datafill/page/mobile/#/backstage/mars-datafill/page/mobile/datafill/collect/usertask"),
		chromedp.WaitVisible(`#root > div.am-pull-to-refresh.am-pull-to-refresh-down > div > div > div:nth-child(2) > div > div.am-accordion-item.am-accordion-item-active > div.am-accordion-header`, chromedp.ByQuery),
		chromedp.FullScreenshot(&pic2, 90),
		// chromedp.Sleep(200*time.Second),
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	netReturn := network.GetRequestPostData("4140")
		// 	log.Println("netReturn:", netReturn)

		// 	return nil
		// }),
	); err != nil {
		log.Fatal(err)
	}

	os.WriteFile("./Screenshots/oaLoginPage.png", pic0, 0o644)
	// os.WriteFile("./Screenshots/oaLoginPageWithAccount.png", pic1, 0o644) // 调试使用，生产环境请注释本行，否则会导致 OA 学号详细泄露
	os.WriteFile("./Screenshots/oaLogined.png", pic2, 0o644)

}

// 判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 监听并拿到 AuthorizationCode
// https://chromedevtools.github.io/devtools-protocol/tot/Network/#event-requestWillBeSent
func getCodeAndDataFill(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {

		// case *network.EventResponseReceived:
		// 	resp := ev.Response
		// 	// if len(resp.Headers) != 0 {
		// 	// 	log.Printf("received urls: %s\n", resp.URL)
		// 	// }
		// 	if resp.URL == "https://oa.hbfu.edu.cn/backstage/mars-datafill/depttask/queryListForPage" {
		// 		log.Printf("queryListForPage URL: %s\n", resp.URL)
		// 		log.Printf("queryListForPage Header: %s\n", resp.Headers)
		// 		log.Printf("queryListForPage RequestHeaders: %s\n", resp.RequestHeaders)
		// 		log.Printf("queryListForPage connectionId: %v\n", resp.ConnectionID)
		// 	}

		case *network.EventRequestWillBeSent:
			// 监听指定方法，找到 queryListForPage 在 CDP 中 Network 的请求，然后解析出 AuthorizationCode ，供后面使用
			if ev.Request.URL == "https://oa.hbfu.edu.cn/backstage/mars-datafill/depttask/queryListForPage" {
				// log.Println("请求完整字段：", ev.Request)                  // 仅调试使用，生产环境务必注释掉本行，防止认证 Code 泄露！
				// log.Println("请求 requestID:", ev.RequestID)              // 仅调试使用，生产环境务必注释掉本行，防止认证 Code 泄露！
				// log.Println("请求 request-postData:", ev.Request.Headers) // 仅调试使用，生产环境务必注释掉本行，防止认证 Code 泄露！
				authorizationHeader, _ := json.Marshal(ev.Request.Headers)
				authorizationHeaderStr := string(authorizationHeader)
				authorizationCodeStrSplit := strings.Split(authorizationHeaderStr, "\"Authorization\":\"")[1]
				AuthorizationCodeStr = strings.Split(authorizationCodeStrSplit, "\",\"Content-Type\"")[0]
				log.Println("登录流程完成!")
				log.Println("AuthorizationCode 获取完成!")
				// log.Println("AuthorizationCode:", AuthorizationCodeStr) // 仅调试使用，生产环境务必注释掉本行，防止认证 Code 泄露！
				dailyDataFill(AuthorizationCodeStr)
			}
		}
	})

}

// 完成当日填写
func dailyDataFill(authorizationCode string) {
	// 取得指定格式的日期用于 OA 查询
	timeLocation, _ := time.LoadLocation("Asia/Shanghai")
	timestamp := time.Now().In(timeLocation)
	nowdate := strftime.Format("%y-%m-%d", timestamp)
	log.Println("今日日期:", nowdate)

	// OA 查询，取得当日未填写的表 ID
	fromID := queryListForPage(authorizationCode, nowdate)
	dailyFill(authorizationCode, fromID, prove, address)

}

// 定义查询结构体
type QueryList struct {
	PageNum   int        `json:"pageNum"`
	PageSize  int        `json:"pageSize"`
	PageParam *pageParam `json:"pageParam"`
}

// 定义查询内容
type pageParam struct {
	Name string `json:"name"`
}

// // 定义 queryListForPage 的查询返回
// type QueryListForPageReturn struct {
// 	Id     string `json:"id"`
// 	UserId string `json:"userId"`
// 	TaskId string `json:"taskId"`
// 	State  int    `json:"state"`
// }

// type QueryListForPageReturn struct {
// 	Id string `json:"id"`
// UserId             string `json:"userId"`
// TaskId             string `json:"taskId"`
// TaskCreateUserName string `json:"taskCreateUserName"`
// TaskName           string `json:"taskName"`
// FormId             string `json:"formId"`
// Content            string `json:"content"`
// }

// 查询当日表格 ID
func queryListForPage(authorizationCode string, nowdate string) (formID string) {
	queryListForPageUrl := "https://oa.hbfu.edu.cn/backstage/mars-datafill/depttask/queryListForPage"
	queryListForPageUrlContentType := "application/json"
	queryParamName := pageParam{
		Name: nowdate,
		// Name: "22-08-29", // 调试用，指定日期
	}
	// 注意赋值方式，否则会报空指针
	queryContent := QueryList{PageNum: 1, PageSize: 10, PageParam: &queryParamName}

	queryData, _ := json.Marshal(queryContent)
	queryParam := bytes.NewBuffer(queryData)

	//构建http请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", queryListForPageUrl, queryParam)
	if err != nil {
		log.Fatal(err)
	}

	// 添加 Header
	request.Header.Add("Authorization", authorizationCode)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
	request.Header.Add("Content-Type", queryListForPageUrlContentType)

	requests, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer requests.Body.Close()
	returnBody, err := io.ReadAll(requests.Body)
	if err != nil {
		log.Fatal(err)
	}
	// returnString := string(returnBody[:])
	// log.Println("returnString:", returnString)
	// var returnData QueryListForPageReturn
	// json.Unmarshal(returnBody, &returnData)
	// log.Println("returnData:", string(returnBody))                             // 仅调试使用，生产环境务必注释掉本行，防止别人拿到表格ID进行填写！
	fillStatus := gjson.Get(string(returnBody[1:len(returnBody)-1]), "state")  // 填写状态
	formID = gjson.Get(string(returnBody[1:len(returnBody)-1]), "id").String() // 表格的 ID ，用于下一步提交
	log.Println("填写状态:", fillStatus)
	// log.Println("formID:", formID) // 仅调试使用，生产环境务必注释掉本行，防止别人拿到表格ID进行填写！
	if fillStatus.Int() != 0 {
		log.Panic("当日数据已经上报，无需再次提交！")
	} else if formID == "" {
		log.Panic("无法找到当日表格 ID ,检查当日打卡是否发放!")
	}
	return formID
}

// 打卡数据结构
type DataFillFull struct {
	Id                string            `json:"id"` // 表单 ID
	*DataFillFormData `json:"formData"` // 表单内容
}

type DataFillFormData struct {
	Address           string `json:"JZDZ"`      // 居住地址（const：address）
	TemperatureNormal string `json:"CWJTWSFZC"` // 晨午检体温是否正常（是）
	Isolateion        string `json:"SFBGL"`     // 是否被隔离（否）
	Prove             string `json:"SFYYSGRQK"` // 是否持有核酸证明（const：prove）
}

// 提交当日打卡数据
func dailyFill(authorizationCode string, formID string, prove bool, address string) {
	log.Println("准备提交当日打卡数据!")

	dailyFillUrl := "https://oa.hbfu.edu.cn/backstage/mars-datafill/formV2/saveandsubmit"
	dailyFillUrlType := "application/json"
	if prove {
		// 持有核酸证明
		HaveProve = "option-2"
	} else {
		// 没有核酸证明
		HaveProve = "option-1"
	}

	// 住址（const：address；温度：正常；隔离：否；核酸证明（const:prove））
	dailyFillFormData := DataFillFormData{Address: address, TemperatureNormal: "是", Isolateion: "option-1", Prove: HaveProve}
	dailyFillContent := DataFillFull{Id: formID, DataFillFormData: &dailyFillFormData}

	dailyFillData, _ := json.Marshal(dailyFillContent)
	dailyFillParam := bytes.NewBuffer(dailyFillData)

	//构建http请求
	client := &http.Client{}
	request, err := http.NewRequest("PUT", dailyFillUrl, dailyFillParam)
	if err != nil {
		log.Fatal(err)
	}

	// 添加 Header
	request.Header.Add("Authorization", authorizationCode)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
	request.Header.Add("Content-Type", dailyFillUrlType)

	requests, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer requests.Body.Close()
	returnBody, err := io.ReadAll(requests.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("提交结果:", string(returnBody))

}

// 判断文件夹是否存在
func hasDir(path string) (bool, error) {
	_, _err := os.Stat(path)
	if _err == nil {
		return true, nil
	}
	if os.IsNotExist(_err) {
		return false, nil
	}
	return false, _err
}

// 创建文件夹
func createDir(path string) {
	_exist, _err := hasDir(path)
	if _err != nil {
		log.Printf("获取截图目录异常 -> %v\n", _err)
		return
	}
	if _exist {
		log.Println("截图目录已已经存在~")
	} else {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Printf("创建截图目录异常 -> %v\n", err)
		} else {
			log.Println("创建截图目录成功!")
		}
	}
}
