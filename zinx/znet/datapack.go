package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"../utils"
	"../ziface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包的头的长度
func (dp *DataPack) GetHeadLen() uint32 {
	//DataLen uint32 (4 bytes) + ID uint342 (4 bytes)
	return 8
}

//封包方法: dataLen|msgID|data
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建存放bytes字节的缓冲
	dataBuf := bytes.NewBuffer([]byte{})

	//将dataLen写进dataBuf中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	//将MsgId写进dataBuf中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//将data写进dataBuf中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

//拆包方法: 将包的Head信息读出来，再根据head中的长度，再进行一次读
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuf := bytes.NewReader(binaryData)

	//只解压head信息，得到dataLen和MsgID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读MsgID
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断dataLen是否超过允许的最大长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}

	//此时msg中并没有data，只有信息头
	return msg, nil
}
