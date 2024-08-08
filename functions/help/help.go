package help

import (
	"log"

	"AIOPrivacyBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendHelpMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	chatID := message.Chat.ID
	messageID := message.MessageID // 获取原始消息ID
	photoPath := "functions/help/help_image.png"
	photoCaption := `欢迎使用 @AIOPrivacyBot ！一个尝试在保护用户隐私的前提下，提供各类功能的类命令行Telegram机器人。

<b>用法</b>

您可以向机器人发送 <code>/help</code> 来获取有关帮助文档。

请通过：<a href="https://github.com/AIOPrivacy/AIOPrivacyBot?tab=readme-ov-file#%E7%94%A8%E6%B3%95">https://github.com/AIOPrivacy/AIOPrivacyBot?tab=readme-ov-file#%E7%94%A8%E6%B3%95</a>访问完整用法支持文档

<blockquote expandable="true">
<b>聊天触发类</b>
<b>/play</b> 动作触发功能
<b>/ask</b> AI提问学术问题
<b>/getid</b> 用户查看ID信息功能
<b>/status</b> 查看系统信息
<b>/admins</b> 召唤所有管理员
<b>/string</b> 字符串编码
<b>/num</b> 数字进制转换
<b>/curconv</b> 货币转换，汇率查询
<b>/color</b> 颜色转换&色卡推荐
支持 <code>RGB</code>、<code>16进制</code>、<code>颜色名称</code>，用法举例：
<code>/color@AIOPrivacyBot #ffffff</code>
<code>/color@AIOPrivacyBot 255 255 255</code>
<code>/color@AIOPrivacyBot Blue</code>

<b>Inline 模式触发类</b>
<b>各类网址的安全过滤/检测</b>
机器人 Inline 模式下运作，您可以这样调用：
<code>@AIOPrivacyBot -check https://www.amazon.com/dp/exampleProduct/ref=sxin_0_pb</code>
<b>内容网站的内容下载存储到 telegraph</b>
机器人 Inline 模式下运作，您可以这样调用：
<code>@AIOPrivacyBot -view https://www.52pojie.cn/thread-143136-1-1.html</code>

<b>其他触发类</b>
<b>回复机器人随机触发 AI聊天触发功能</b>
70% 概率的出现 <b>笨笨的猫娘</b> AI 玩耍！
</blockquote>
↑ 点击展开详细命令说明

⚠ 机器人正在测试中，如果遇到bug请及时提出！

欢迎加入用户交流群：<a href="https://t.me/AIOPrivacy">https://t.me/AIOPrivacy</a>`

	// 使用 utils.SendPhotoWithCaption 发送带有文字的图片并回复到用户
	err := utils.SendPhotoWithCaption(chatID, messageID, photoPath, photoCaption, bot)
	if err != nil {
		log.Printf("Error sending help image: %v", err)
	}
}
