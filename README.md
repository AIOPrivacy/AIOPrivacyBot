<a id="readme-top"></a>

<div align="center">

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![GPL-3.0 License][license-shield]][license-url]

</div>

<br />

<div align="center">
  <a href="https://github.com/iuu6/AIOPrivacyBot">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a>
  <h3 align="center">AIOPrivacyBot</h3>

  <p align="center">
    一个尝试在保护用户隐私的前提下，提供各类功能的类命令行Telegram机器人。
    <br />
    <a href="https://github.com/iuu6/AIOPrivacyBot"><strong>帮助文档 »</strong></a>
    <br />
    <br />
    <a href="https://t.me/AIOPrivacyBot">查看演示</a>
    ·
    <a href="https://github.com/iuu6/AIOPrivacyBot/issues">报告错误/提出更多内容</a>
    ·
    <a href="https://t.me/AIOPrivacy">加入用户交流群组</a>
  </p>
</div>




<details>
  <summary>项目目录</summary>
  <ol>
    <li>
      <a href="#关于此项目">关于此项目</a>
      <ul>
        <li><a href="#由谁构建">由谁构建</a></li>
      </ul>
    </li>
    <li>
      <a href="#快速开始">快速开始</a>
      <ul>
        <li><a href="#先决条件">先决条件</a></li>
        <li><a href="#安装">安装</a></li>
      </ul>
    </li>
    <li><a href="#用法">用法</a></li>
    <li><a href="#功能路线">功能路线</a></li>
    <li><a href="#贡献指南">贡献指南</a></li>
  </ol>
</details>



# 关于此项目

![Product Name Screen Shot][product-screenshot]

本项目一直开启Telegram Bot自带的Privacy Mode，不赋予管理员权限的情况下，只可以读取@机器人的指令和回复机器人的消息，从而极大保护用户隐私。

各类需要处理的消息也**将**会采用签名的方法保证消息未篡改，并对部分敏感数据进行哈希加盐。

<p align="right">(<a href="#readme-top">back to top</a>)</p>



## 由谁构建

本项目主要由Golang语言构建

<p align="right">(<a href="#readme-top">back to top</a>)</p>



# 快速开始

## 先决条件

有一台可以正常运行二进制文件的电脑

## 安装

直接运行构建后的二进制文件即可。

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# 用法

您可以向机器人发送`/help`来获取有关帮助文档。

## 聊天触发类

### /play 动作触发功能

以下演示都以**A回复B**模拟！

##### 主动模式

`/play@AIOPrivacyBot -t xxxxx` 可以成功触发`A xxxxx了 B！`

`/play@AIOPrivacyBot -t xxxxx yyyyy` 可以成功触发`A xxxxx B yyyyy`

##### 被动模式

`/play@AIOPrivacyBot -p xxxxx` 可以成功触发`A 被 B xxxxx了！`

`/play@AIOPrivacyBot -p xxxxx yyyyy` 可以成功触发`B xxxxx A yyyyy`

##### 备注
注意：可以使用英文'或"包括发送内容来高于空格优先级，例如`/play@AIOPrivacyBot -p "xx xxx" "yy yy y"`

### /ask AI提问学术问题

在私聊或群聊中均可使用，发送`/ask@AIOPrivacyBot`即可触发，调用gpt-4o-mini来解决较为严谨的学术问题

### /getid 用户查看ID信息功能

您可以在私聊或群聊中发送`/getid@AIOPrivacyBot`或`/getid`，来获取自己和群组详细的Telegram ID等信息

### /status 查看系统信息

您可以在私聊或群聊中发送`/status@AIOPrivacyBot`或`/status`，来查看机器人和系统的运行状态

### /admins 召唤所有管理员

`/admins@AIOPrivacyBot`即可召唤本群所有管理员（危险功能，需要确认后才会@管理员）

### /string 字符串编码

`/string@AIOPrivacyBot -url xxx`即可进行字符串转换
包括多个参数，可以使用`-all`查看

### /num 数字进制转换

`/string@AIOPrivacyBot`可以进行进制转换

可以输入整数，也可以输入0xfff（十六进制）

```
b 表示 bin
o 表示 oct
x 表示 hex
```

### /curconv 货币转换，汇率查询

`/curconv@AIOPrivacyBot`可以进行货币转换，汇率查询

### /color 颜色转换&色卡推荐

`/color@AIOPrivacyBot`获取颜色转换&色卡推荐，支持`RGB` `16进制` `颜色名称`

用法举例

```
/color@AIOPrivacyBot #ffffff
/color@AIOPrivacyBot 255 255 255
/color@AIOPrivacyBot Blue
```

## Inline 模式触发类

### 各类网址的安全过滤/检测

机器人Inline模式下运作，您可以这样调用

```
@AIOPrivacyBot -check https://www.amazon.com/dp/exampleProduct/ref=sxin_0_pb?__mk_de_DE=%C3%85M%C3%85%C5%BD%C3%95%C3%91&keywords=tea&pd_rd_i=exampleProduct&pd_rd_r=8d39e4cd-1e4f-43db-b6e7-72e969a84aa5&pd_rd_w=1pcKM&pd_rd_wg=hYrNl&pf_rd_p=50bbfd25-5ef7-41a2-68d6-74d854b30e30&pf_rd_r=0GMWD0YYKA7X
```

### 内容网站的内容下载存储到telegraph

机器人Inline模式下运作，您可以这样调用

```
@AIOPrivacyBot -view https://www.52pojie.cn/thread-143136-1-1.html
```

## 其他触发类

### 回复机器人随机触发 AI聊天触发功能

70%概率的出现**笨笨的猫娘**AI玩耍！



<p align="right">(<a href="#readme-top">back to top</a>)</p>

# 功能路线

- [x] 指令化重构

- [x] /play 动作触发功能

- [x] 回复随机触发 AI聊天触发功能

- [x] /ask AI提问学术问题

- [x] /help 帮助中心

- [x] /getid 用户查看ID信息功能

- [x] /status 查看系统信息

- [x] /admins 召唤所有管理员

- [x] /num 数字进制转换

- [x] /string 字符串编码

- [x] /curconv 货币转换，汇率查询

- [x] /color 颜色转换&色卡推荐

- [x] 各类网址的安全检测

- [x] 各类网址的安全过滤

- [x] CSDN/吾爱破解/知乎/……等等诸多内容网站的内容下载存储到telegraph，避免隐私窃取

- [ ] 功能开关支持

- [ ] 多语言支持

- [ ] 消息发送安全性确认功能

- [ ] 自检程序有无被修改功能

  **以下功能只适用于机器人有群组管理员的情况**

- [ ] 入群欢迎功能

- [ ] 入群验证功能

- [ ] 自动移除所有非管理员用户

- [ ] 自动取消频道消息的置顶

- [ ] fban - 封禁联盟功能

<p align="right">(<a href="#readme-top">back to top</a>)</p>




# 贡献指南

正是贡献让开源社区成为了学习、启发和创造的绝佳场所。我们**非常感谢**您的任何贡献。

如果您有改进建议，请分叉该仓库并创建拉取请求。您也可以直接打开一个带有标签“增强”的问题。

*等待完善*

<p align="right">(<a href="#readme-top">back to top</a>)</p>


[contributors-shield]: https://img.shields.io/github/contributors/iuu6/AIOPrivacyBot.svg?style=for-the-badge
[contributors-url]: https://github.com/iuu6/AIOPrivacyBot/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/iuu6/AIOPrivacyBot.svg?style=for-the-badge
[forks-url]: https://github.com/iuu6/AIOPrivacyBot/network/members
[stars-shield]: https://img.shields.io/github/stars/iuu6/AIOPrivacyBot.svg?style=for-the-badge
[stars-url]: https://github.com/iuu6/AIOPrivacyBot/stargazers
[issues-shield]: https://img.shields.io/github/issues/iuu6/AIOPrivacyBot.svg?style=for-the-badge
[issues-url]: https://github.com/iuu6/AIOPrivacyBot/issues
[license-shield]: https://img.shields.io/github/license/iuu6/AIOPrivacyBot.svg?style=for-the-badge
[license-url]: https://github.com/iuu6/AIOPrivacyBot/blob/master/LICENSE
[product-screenshot]: images/screenshot.png