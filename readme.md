# go dlt645-2007/1997




<img src="https://img.shields.io/github/stars/cqqqq/godlt645?style=social"/>

用go语言实现的dlt645解析
```shell
    go get github.com/cqqqq/godlt645
```

1. 读请求
```go
c := godlt645.NewClient(godlt645.NewRTUClientProvider(godlt645.WithEnableLogger(), godlt645.WithSerialConfig(serial.Config{
    Address:  "/dev/ttyUSB3",
    BaudRate: 19200,
    DataBits: 8,
    StopBits: 1,
    Parity:   "E",
    Timeout:  time.Second * 8,
})))

for {
    time.Sleep(time.Second)
    pr, err := c.Read(godlt645.NewAddress("3a2107000481", godlt645.LittleEndian), 0x00_01_00_00)
    if err != nil {
        log.Print(err.Error())
    } else {
        println(pr.GetValue())
}

}

```
2. protocal_test.go 中可查看更多案例
