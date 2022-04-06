# go-third-login
## 第一次提交
封装了微信小程序第三方登录，只需传入wx.login得到的code，以及微信自带的getPhoneNumer的e.detail.code，并且在application.yaml下配置wechat.appid以及wechat.secret,返回结果为手机号，openid，error，openid可以作为userid

注：目前仅为方便自己以后其他业务方面调用，但也欢迎大家测试是否存在安全漏洞，提出issue.
