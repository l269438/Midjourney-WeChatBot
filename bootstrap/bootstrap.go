package bootstrap

import (
	"fmt"
	"github.com/869413421/wechatbot/handlers"
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
	"log"
	"net/url"
	"os"
	"strings"
)

type InputData struct {
	//ImageURL  string `json:"imageUrl"`
	ImageURL  string `json:"uri"`
	MessageID string `json:"messageId"`
	State     string `json:"state"`
	MsgHash   string `json:"msgHash"`
	Prompt    string `json:"prompt"`
	PromptEn  string `json:"promptEn"`
	Id        string `json:"id"`
	Action    string `json:"action"`
}

type ResponseData struct {
	ImageURL string `json:"uri"`
	Id       string `json:"id"`
	State    string `json:"state"`
	MsgHash  string `json:"hash"`
	Prompt   string `json:"prompt"`
	PromptEn string `json:"promptEn"`
	Action   string `json:"action"`
}

func crontab(self *openwechat.Self, groups openwechat.Groups, body ResponseData) {
	fmt.Println("进入回调")
	fmt.Println(body)
	promptEn := body.PromptEn
	prompt := body.Prompt

	imgUrl := body.ImageURL
	state := body.State
	action := body.Action
	id := body.Id
	parts := strings.Split(state, ":")
	atText := "@" + parts[1]

	name := groups.GetByNickName(parts[0])
	parsedURL, err := url.Parse(imgUrl)
	if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
		log.Printf("ImageUrl :%v \n", imgUrl)
		tmpImageFile, err := handlers.DownloadImage(imgUrl)
		defer tmpImageFile.Close()
		tmpImageFile.Seek(0, 0) // 将文件指针重置到文件开头
		if err != nil {
			log.Printf("download image error: %v \n", err)
			//msg.ReplyText(data)
		} else if name != nil {
			if action == "UPSCALE" {
				self.SendImageToGroup(name, tmpImageFile)
				result := "✅绘制成功\n" +
					"\n"
				self.SendTextToGroup(name, atText+result)
			} else {
				self.SendImageToGroup(name, tmpImageFile)
				result := "✅绘制成功\n" +
					"📎任务ID: " + id + "\n" +
					"\n" +
					"🙋🏻 Prompt: " + prompt + "\n" +
					"\n" +
					"✏️ PromptEn: " + promptEn + "\n" +
					"\n" +
					"🪄 放大：这里有四幅草图，请用 U+编号来告诉我您喜欢哪一张。例如，第一张为U1。我将会根据您的选择画出更精美的版本。" +
					"\n" +
					"🪄 变换：如果您对所有的草图都不太满意，但是对其中某一张构图还可以，可以用 V+编号来告诉我，我会画出类似的四幅草图供您选择" +
					"\n" +
					"✏ 具体操作：[ex 编号,操作]，比如 ex 0234495019546343,U1"
				//self.SendTextToGroup(name, atText+" 您的图片已生成标识符为："+msgHash)
				self.SendTextToGroup(name, atText+result)
			}

			defer os.Remove(tmpImageFile.Name())
		}
	} else {
	}

}

func Run() {
	//bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	// 注册消息处理函数
	bot.MessageHandler = handlers.Handler
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
	// 执行热登录
	err := bot.HotLogin(reloadStorage)
	if err != nil {
		if err = bot.Login(); err != nil {
			log.Printf("login error: %v \n", err)
			return
		}
	}

	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(self)

	groups, err := self.Groups()

	// 获取所有的群组
	if err != nil {
		log.Println(err)
	}
	fmt.Println(groups, err)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	r := gin.Default()

	r.POST("/mj/v3/webhook", func(c *gin.Context) {
		var body ResponseData
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// 在闭包中访问 self 和 groups 变量
		go crontab(self, groups, body)

		c.JSON(200, gin.H{"message": "cron task started"})
	})

	go func() {
		if err := r.Run(":9095"); err != nil {
			panic(err)
		}
	}()

	bot.Block()
}
