package godlt645

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Exception 错误
type Exception struct {
	Err uint16
}

func (e *Exception) Error() string {

	return fmt.Sprintf("645 err %d", e.Err)
}

func (e Exception) Encode(buffer *bytes.Buffer) error {
	return binary.Write(buffer, binary.LittleEndian, e.Err+0x33)
}

// GetLen 错误响应报文长度2
func (e Exception) GetLen() byte {
	return 2
}
