## 编译

```
go build -ldflags="-linkmode internal -H windowsgui"
```

## 技术实现
- golang
- lxn/walk
- SSH

## 功能
- 扫描指定网段下的THOR矿机
- 矿机重启
- 修改矿工配置并自动重启
- 查询矿机运行日志

## 功能预览

矿机详情
![image](https://github.com/langubtc/ThorClient/blob/master/img/run.png)

矿机重启
![image](https://github.com/langubtc/ThorClient/blob/master/img/view.png)