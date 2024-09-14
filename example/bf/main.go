package main

import (
	"bytes"
	"flag"
	"github.com/cqqqq/godlt645"
	"github.com/goburrow/serial"
	"io"
	"log"
	"time"
)

var _ godlt645.PrefixHandler = (*Handler)(nil)

type Handler struct {
}

func (h Handler) EncodePrefix(buffer *bytes.Buffer) error {
	buffer.Write([]byte{0xfe, 0xfe, 0xfe, 0xfe})
	return nil
}

func (h Handler) DecodePrefix(reader io.Reader) ([]byte, error) {
	fe := make([]byte, 4)
	_, err := io.ReadAtLeast(reader, fe, 4)
	if err != nil {
		return nil, err
	}
	return fe, err
}

// 百富电表
func main() {
	var b int
	var code int
	flag.IntVar(&b, "b", 2400, "波特率")
	flag.IntVar(&code, "c", 0x00_03_00_00, "波特率")
	flag.Parse()
	p := godlt645.NewRTUClientProvider(godlt645.WithSerialConfig(serial.Config{Address: "COM2", BaudRate: b, DataBits: 8, StopBits: 1, Parity: "E", Timeout: time.Second * 30}), godlt645.WithEnableLogger(), godlt645.WithPrefixHandler(&Handler{}))
	c := godlt645.NewClient(p)
	err := c.Connect()
	if err != nil {
		panic(err)
	}
	defer c.Close()
	for {
		time.Sleep(50 * time.Millisecond)
		read, _, err := c.Read(godlt645.NewAddress("000703200136", godlt645.LittleEndian), 0x00_01_00_00)
		if err == nil {
			log.Printf("rec %s", read.GetValue())
		}

	}
}
