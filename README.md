# HNU官微定时打卡程序

下载到本地 git clone https://github.com/SheyQ/Clock-in.git

打开文件 cd Clock-in

编辑config.json
```
{
    "Time": "00:30:00",     // 打卡的时间
    "Code": "2016260103xx", // 学号
    "Password": "xxxxxx", // 门户密码
    "RealProvince": "xx省", // 省
    "RealCity": "xx市",       //市
    "RealCounty": "xx区",     //区
    "RealAddress": "xx街道404" //详细地址
}
```

测试运行: go run clock.go test

后台运行: nohup go run clock.go &