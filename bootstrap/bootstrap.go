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
	ImageURL  string `json:"imageUrl"`
	MessageID string `json:"messageId"`
	State     string `json:"state"`
	MsgHash   string `json:"msgHash"`
}

type Request struct {
	Action      string `json:"action"`
	ID          string `json:"id"`
	Prompt      string `json:"prompt"`
	Description string `json:"description"`
	State       string `json:"state"`
	SubmitTime  int64  `json:"submitTime"`
	FinishTime  *int64 `json:"finishTime"`
	ImageURL    string `json:"imageUrl"`
	Status      string `json:"status"`
}

func crontab(self *openwechat.Self, groups openwechat.Groups, body Request) {
	fmt.Println("进入回调")
	fmt.Println(body)
	fmt.Println("返回的ImageUrl:" + body.ImageURL)

	imgUrl := body.ImageURL
	state := body.State
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
			self.SendImageToGroup(name, tmpImageFile)
			self.SendTextToGroup(name, atText+" 您的图片已生成标识符为："+body.ID)
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
		var body Request
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
