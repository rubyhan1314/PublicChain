package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
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
