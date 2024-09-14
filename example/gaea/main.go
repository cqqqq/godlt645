package main

import (
	"bytes"
	"flag"
	"github.com/cqqqq/godlt645"
	"github.com/goburrow/serial"
	"log"
	"sync"
	"time"
)

var device sync.Map
var mu sync.Mutex

// 特殊电表解析 同步
func main() {
	var b int
	var code int
	//1200、2400、4800、9600
	flag.IntVar(&b, "b", 19200, "波特率")
	flag.IntVar(&code, "c", 0x02_04_00_00, "特征码")
	flag.Parse()
	p := godlt645.NewRTUClientProvider(
		godlt645.WithSerialConfig(serial.Config{
			Address:  "/dev/ttyUSB3",
			BaudRate: b,
			DataBits: 8,
			StopBits: 1,
			Parity:   "E",
			Timeout:  time.Second * 30}),
		godlt645.WithEnableLogger())
	c := godlt645.NewClient(p)
	err := c.Connect()
	if err != nil {
		panic(err)
	}
	if err != nil {
		return
	}
	defer c.Close()

	forceOnline(c)

	go func() {
		time.Sleep(1 * time.Minute)
		forceOnline(c)
	}()

	for {
		//如果对扫描速度不是要求很高 需要重新扫描
		time.Sleep(50 * time.Millisecond)

		device.Range(func(key, value interface{}) bool {
			func() {
				mu.Lock()
				defer mu.Unlock()
				read, _, err := c.Read(godlt645.NewAddress(key.(string), godlt645.LittleEndian), int32(code))
				if err == nil {
					log.Printf("address %s :rec %f", key.(string), read.GetFloat64Value())
				}
			}()
			return false
		})
	}

}
func forceOnline(c godlt645.Client) {
	c.Broadcast(godlt645.NullData{}, *godlt645.NewControlValue(0x0a))
	mu.Lock()
	defer mu.Unlock()
	for {

		frame, err := c.ReadRawFrame()
		if err != nil {
			log.Printf(err.Error())
			return
		}
		if frame != nil && len(frame) > 10 {
			p, err := godlt645.Decode(bytes.NewBuffer(frame))
			if err != nil {
				log.Printf(err.Error())
				return
			}
			device.Store(p.Address.GetStrAddress(godlt645.LittleEndian), nil)
		}
	}
}
