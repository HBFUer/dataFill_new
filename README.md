# dataFill_new
河北金融学院每日自动健康打卡（新版）

适配了最新微信小程序打卡方式，适用于自 2022-06-22 后的每日健康打卡

## 运行截图

如果程序顺利运行，将会在运行中截取两张图片放在 `./Screenshots` 文件夹下，分别是：

（tips：点击箭头可以展开查看图片）

<details>   <summary>oaLoginPage.png OA 登录页截图</summary>   <p><img src=./images/oaLoginPage.png/></p>    <span>OA 登录页截图</span></details>

<details>   <summary>oaLogined.png 健康打卡页截图</summary>   <p><img src=./images/oaLogined.png/></p>    <span>健康打卡页截图</span></details>

运行结果会是以下情况之一：

<details>   <summary>当日未打卡，完成打卡</summary>   <p><img src=./images/successful.png/></p></details>

<details>   <summary>当日打卡未发放（一般是没有到中午12点）</summary>   <p><img src=./images/fail0.png/></p></details>

<details>   <summary>当日打卡已经完成，不需要重复提交</summary>   <p><img src=./images/fail1.png/></p></details>

## 运行方法

我们提供了 `GitHub Action`（推荐） 和 `单文件` 两种运行方式

### GitHub Action（推荐）

> GitHub Action 是 GitHub （2018年被微软收购）提供的服务，可以利用微软的服务器运行自己的程序，GitHub 为每个个人用户提供了每月 3000 分钟的免费运行额度，在此感谢 GitHub 和微软！

> 使用这种方式，你在设定好之后每日即可自动打卡，不需要自己准备设备，不需要人工干预！

> 请在理解单文件运行方式的前提下使用 GitHub Action 

1. 首先你需要拥有一个 GitHub 账号，如果没有，请点击[这里](https://github.com/)然后点击右上角的 `Sign up` 注册一个，可以看这篇教程：https://blog.csdn.net/KKKKKKKMAx/article/details/125394719
2. **fork** 本项目到自己的仓库，**fork** 按钮在这个页面的右上角**（边上还有一个 `Star` 的按钮，希望你一并点亮支持一下作者）**，如果不会的话，参考：https://blog.csdn.net/weixin_52634719/article/details/122504390
3. 点击Actions选项卡，点击`I understand my workflows, go ahead and enable them`
4. 在 **fork 下来的自己的仓库中**到 **Settings-Secrets** 设置密钥（把密码等信息放在这里别人是看不到的），如图设置四个 `Repository secrets` ，如果不会的话，请参考：https://zhuanlan.zhihu.com/p/516839666 中“设置Github secrets 环境变量”一节

![](./images/secretSetting.png)

字段说明：

| 字段       | 说明                                                         | 举例                        | 备注                                                       |
| ---------- | ------------------------------------------------------------ | --------------------------- | ---------------------------------------------------------- |
| oaUsername | 河北金融学院OA用户名（也就是学工号）                         |                             | 填写时复制本字段，会自动转大写，下同                       |
| oaPassword | 河北金融学院OA密码                                           | 123456                      |                                                            |
| address    | 自己的居住地址，请填写URL编码（encodeURL编码）的地址，直接填写中文会导致上传的地址数据乱码 | %E4%B8%8B%E5%8C%97%E6%B3%BD | URL在线编码网站参考：https://www.bejson.com/enc/urlencode/ |
| prove      | 是否持有核酸证明是（true）否（false）                        | true                        |                                                            |

> tips:使用 GitHub Action 时，`prove` 是必选参数

完成，程序将在北京时间（UTC+8）的16时07分进行自动打卡操作，如果需要修改这个时间，请修改`./github/workflows/dailyFill.yml` 文件中的 `cron` 字段

> tips:使用 GitHub Action 时，根据当时的排队情况，程序运行的时间可能会稍有延迟（一般在半个小时以内），所以可能在你设定时间之后的一段时间才能完成打卡

### 单文件

#### 下载

到 [Release（点击链接前往）](https://github.com/HBFUer/dataFill_new/releases)  ，根据你使用的硬件，下载最新版的文件，比如你是 Windows 的电脑就下载 `datafill_new_windows_amd64.exe` ，如果是 Mac 电脑，就下载 `datafill_new_darwin_amd64`（Intel芯片）或者 `datafill_new_darwin_arm64`（M1芯片），如果是在路由器（OpenWRT）上跑的话，一般是 mips 架构的，就是 `datafill_new_linux_mips`

#### 运行

【前置】你需要安装 Chrome 浏览器或者 Headless Chrome（没有 GUI 的话）

直接执行上一步下载的文件就能看到相关提示

![](./images/run0.png)

比如你是法外狂徒`张三`，你使用`Windows`电脑，你的学号（或者是工号）是`114514`，你的[OA](https://oa.hbfu.edu.cn/backstage/cas/login)密码是`1919810`，你居住在`河北省保定市莲池区下北泽街道3188号河北金融学院`，拥有核酸检测证明，那么你应当在程序目录下打开 PowerShell 或者 cmd，输入如下参数完成当日打卡

```powershell
.\datafill_new_windows_amd64.exe -oaUsername=114514 -oaPassword=1919810 -address=河北省/保定市/莲池区/下北泽街道/3188号河北金融学院 -prove=true
```

> tips:省/市/区/街道（乡）的具体填写请参考 https://oa.hbfu.edu.cn/datafill/collect/usertask 行政规划之间使用英文`/`分隔，不需要加空格

如果`张三`没有核酸检测证明，则输入

```powershell
.\datafill_new_windows_amd64.exe -oaUsername=114514 -oaPassword=1919810 -address=河北省/保定市/莲池区/下北泽街道/3188号河北金融学院 -prove=false
```

> tips:此时 -prove 参数可以省略

## 注意事项

- 本程序将自动完成每日健康打卡，你需要对你上报的数据负责！程序仅负责调用接口上报数据！
- 程序仅供学习探讨Go语言编程，对使用本程序造成的一切后果作者均不负责！
- 程序不存储用户账户密码，请妥善保管好相关信息！
- 程序不对接口变动后可能产生的异常负责，请关注接口信息！
- 运行程序则代表已知晓并同意以上规则！

## 反馈&Bug上报

请在 https://github.com/HBFUer/dataFill_new 仓库下发起 Issues

留下：

- 你使用的程序版本
- 你的系统版本，设备信息
- 你的问题，使用上的疑问也可以

# FAQ

1. 使用 GitHub Action 运行，有什么办法可以暂停每日运行？

> A:你可以在 Action 里禁用掉这个 Workflow ，选择名为“每日健康打卡” 的 Workflow ，然后点击 Disable ，需要运行的时候在 Enable ，如果 fork 之后没有正常按时运行，也请检查这里的情况
>
> ![](./images/disableAction.png)

## 更新历史

### 1.01

- 尝试解决 GitHub Action Secret 的中文乱码问题
