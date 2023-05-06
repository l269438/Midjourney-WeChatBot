package handlers

import (
	"bufio"
	"fmt"
	"github.com/869413421/wechatbot/gtp"
	"github.com/eatmoreapple/openwechat"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var badWords []string

// 初始化敏感词列表
func init() {
	var err error
	badWords, err = loadBadWordsFromFile("./profanities.txt")
	if err != nil {
		log.Fatalf("Failed to load bad words: %v", err)
	}
}

type RequestLimiter struct {
	sync.Mutex
	requestCount int
	resetTime    time.Time
}

func (rl *RequestLimiter) CanRequest() bool {
	rl.Lock()
	defer rl.Unlock()

	now := time.Now()
	if now.After(rl.resetTime) {
		rl.requestCount = 0
		rl.resetTime = now.Add(1 * time.Minute)
	}

	if rl.requestCount >= 3 {
		return false
	}

	rl.requestCount++
	return true
}

var limiter = &RequestLimiter{}

func loadBadWordsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var badWords []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			badWords = append(badWords, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return badWords, nil
}

func containsBadWords(text string, badWords []string) bool {
	text = strings.ToLower(text)
	for _, badWord := range badWords {
		badWord = strings.ToLower(badWord)
		if strings.Contains(text, badWord) {
			return true
		}
	}
	return false
}

var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)

// GroupMessageHandler 群消息处理
type GroupMessageHandler struct {
}

// handle 处理消息
func (g *GroupMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyImg(msg)
	}
	return nil
}

// NewGroupMessageHandler 创建群消息处理器
func NewGroupMessageHandler() MessageHandlerInterface {
	return &GroupMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *GroupMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收群消息
	sender, err := msg.Sender()
	group := openwechat.Group{sender}
	log.Printf("Received Group %v Text Msg : %v", group.NickName, msg.Content)

	// 不是@的不处理
	if !msg.IsAt() {
		return nil
	}

	// 替换掉@文本，然后向GPT发起请求
	replaceText := "@" + sender.Self.NickName
	requestText := strings.TrimSpace(strings.ReplaceAll(msg.Content, replaceText, ""))
	reply, err := gtp.Completions(requestText)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		msg.ReplyText("机器人神了，我一会发现了就去修。")
		return err
	}
	if reply == "" {
		return nil
	}

	// 获取@我的用户
	groupSender, err := msg.SenderInGroup()
	if err != nil {
		log.Printf("get sender in group error :%v \n", err)
		return err
	}

	// 回复@我的用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	atText := "@" + groupSender.NickName
	replyText := atText + reply
	_, err = msg.ReplyText(replyText)
	if err != nil {
		log.Printf("response group error: %v \n", err)
	}
	return err
}
func (g *GroupMessageHandler) ReplyImg(msg *openwechat.Message) error {
	if !msg.IsAt() {
		return nil
	}
	if strings.Contains(msg.Content, "help") && msg.IsAt() {
		result := "欢迎使用MJ机器人\n" +
			"------------------------------\n" +
			"🎨 生成图片命令\n" +
			"输入: mj prompt\n" +
			"prompt 即你向mj提的绘画需求\n" +
			"------------------------------\n" +
			"🌈 变换图片命令\n" +
			"输入: ex 标识符 U1\n" +
			"输入: ex 3939314233586510,V1\n" +
			"3939314233586510代表任务ID，U代表放大，V代表细致变化，1代表第1张图 需要逗号隔开\n" +
			"------------------------------\n" +
			"📕 附加参数 \n" +
			"1.解释：附加参数指的是在prompt后携带的参数，可以使你的绘画更加别具一格\n" +
			"· 输入 mj prompt --v 5 --ar 16:9\n" +
			"2.使用：需要使用--key value ，key和value之间需要空格隔开，每个附加参数之间也需要空格隔开\n" +
			"------------------------------\n" +
			"📗 附加参数列表\n" +
			"1.(--version) 或 (--v) 《版本》 参数 1，2，3，4，5 ，不可与niji同用\n" +
			"2.(--niji)《卡通版本》 参数 空或 5 默认空，不可与版本同用\n" +
			"3.(--aspect) 或 (--ar) 《横纵比》 参数 n:n ，默认1:1\n" +
			"4.(--chaos) 或 (--c) 《噪点》参数 0-100 默认0\n" +
			"5.(--quality) 或 (--q) 《清晰度》参数 .25 .5 1 2 分别代表，一般，清晰，高清，超高清，默认1\n" +
			"6.(--style) 《风格》参数 4a,4b,4c (v4)版本可用，参数 expressive,cute (niji5)版本可用\n" +
			"7.(--stylize) 或 (--s)) 《风格化》参数 1-1000 v3 625-60000\n" +
			"8.(--seed) 《种子》参数 0-4294967295 可自定义一个数值配合(sameseed)使用\n" +
			"9.(--sameseed) 《相同种子》参数 0-4294967295 可自定义一个数值配合(seed)使用\n" +
			"10.(--tile) 《重复模式》参数 空"
		msg.ReplyText(result)
		return nil
	}
	if !limiter.CanRequest() {
		msg.ReplyText("请求太快了，请在一分钟后再试。")
		return nil
	}
	if containsBadWords(msg.Content, badWords) {
		msg.ReplyText("您的消息中包含敏感词，请修改后再发送。")
		return nil
	}
	maxInt := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(5)
	time.Sleep(time.Duration(maxInt+1) * time.Second)
	// 接收群消息
	sender, err := msg.Sender()
	group := openwechat.Group{sender}
	log.Printf("Received Group %v Text Msg : %v", group.NickName, msg.Content)
	groupSender, err := msg.SenderInGroup()
	atText := "@" + groupSender.NickName

	state := group.NickName + ":" + groupSender.NickName
	if strings.Contains(msg.Content, "mj") {
		replaceText := "@midjourney"
		requestText := strings.TrimSpace(strings.ReplaceAll(msg.Content, replaceText, ""))
		requestText = strings.TrimSpace(strings.Replace(requestText, "mj", "", 1))
		//messageId, err := gtp.GetMessageId(requestText)

		messageId, err := gtp.GetMessageId(requestText, state, "IMAGINE")
		fmt.Println("请求返回的" + messageId)
		if err != nil {
			log.Printf("gtp request error: %v \n", err)
			msg.ReplyText("超时了 请稍后再试。")
			return err
		}
		if messageId != "" {
			fmt.Println("群名称" + group.NickName)
			fmt.Println("用户名称" + groupSender.NickName)
			msg.ReplyText(atText + "正在生成图片，请稍等...")
		}
	} else if strings.Contains(msg.Content, "ex") {
		replaceText := "@" + sender.Self.NickName
		requestText := strings.TrimSpace(strings.ReplaceAll(msg.Content, replaceText, ""))
		requestText = strings.TrimSpace(strings.Replace(requestText, "ex", "", 1))

		dataParts := strings.Split(requestText, ",")

		if len(dataParts) >= 2 {
			var buttonMessageId = strings.TrimSpace(dataParts[0])
			var button = strings.TrimSpace(dataParts[1])

			fmt.Printf("Button Message ID: %s\n", buttonMessageId)
			fmt.Printf("Button: %s\n", button)

			action, _, err := buttonAction(button)
			if action == "error" {
				msg.ReplyText("传入标识符有误")
			}

			messageId, err := gtp.GetEx(state, action, button, buttonMessageId)
			fmt.Println("请求返回的" + messageId)
			if err != nil {
				log.Printf("gtp request error: %v \n", err)
				msg.ReplyText("超时了 请稍后再试。")
				return err
			}
			if messageId != "" {
				fmt.Println("群名称" + group.NickName)
				fmt.Println("用户名称" + groupSender.NickName)
				msg.ReplyText(atText + "正在生成图片，请稍等...")
			}
		} else {
			fmt.Println("Invalid input format.")
		}

	} else {
		g.ReplyText(msg)
	}
	return err
}
func buttonAction(button string) (string, int64, error) {
	validButtons := []string{"V1", "V2", "V3", "V4", "U1", "U2", "U3", "U4"}

	// Check if the button is in the validButtons array, ignoring case
	isButtonValid := false
	for _, validButton := range validButtons {
		if strings.EqualFold(button, validButton) {
			isButtonValid = true
			break
		}
	}

	if !isButtonValid {
		return "", 0, fmt.Errorf("error")
	}

	// Check the button value and return the corresponding output
	var actionType string
	var index int64
	if strings.HasPrefix(strings.ToUpper(button), "V") {
		actionType = "VARIATION"
		indexString := strings.TrimPrefix(strings.ToUpper(button), "V")
		indexValue, err := strconv.ParseInt(indexString, 10, 64)
		if err != nil {
			return "", 0, err
		}
		index = indexValue
	} else if strings.HasPrefix(strings.ToUpper(button), "U") {
		actionType = "UPSCALE"
		indexString := strings.TrimPrefix(strings.ToUpper(button), "U")
		indexValue, err := strconv.ParseInt(indexString, 10, 64)
		if err != nil {
			return "", 0, err
		}
		index = indexValue
	}

	return actionType, index, nil
}

func DownloadImage(imageURL string) (*os.File, error) {
	// 发起 GET 请求下载图像
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 获取文件名
	urlPath := resp.Request.URL.Path
	filename := filepath.Base(urlPath)

	// 创建临时文件
	tmpFile, err := ioutil.TempFile("", filename+"_*"+".jpg")
	if err != nil {
		return nil, err
	}

	// 将下载的内容写入临时文件
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		// 关闭并删除临时文件
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, err
	}

	// 返回临时文件的句柄
	return tmpFile, nil
}
