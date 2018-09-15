package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
	"encoding/gob"
	"fmt"
)

/*
将一个int64的整数：转为二进制后，每8bit一个byte。转为[]byte
 */
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	//将二进制数据写入w
	//func Write(w io.Writer, order ByteOrder, data interface{}) error
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	//转为[]byte并返回
	return buff.Bytes()
}

/*
Json字符串转为[] string数组
 */
func JSONToArray (jsonString string) [] string{
	var sArr [] string
	if err := json.Unmarshal([]byte(jsonString),&sArr);err != nil{
		log.Panic(err)
	}
	return sArr
}


//字节数组反转
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}


//version 转字节数组
func commandToBytes(command string) []byte {
	var bytes [COMMANDLENGTH]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}



//字节数组转command
func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}


// 将结构体序列化成字节数组
func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
