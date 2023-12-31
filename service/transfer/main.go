package main

import (
	"bufio"
	"encoding/json"
	"filestore-server/config"
	"filestore-server/db"
	"filestore-server/mq"
	"filestore-server/store/oss"
	"log"
	"os"
)

func ProcessTransfer(msg []byte) bool {
	//1.解析msg

	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	//2.根据临时存储文件路径，创建文件句柄
	filed, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	//3.通过文件句柄将文件内容读出来写入oss
	oss.Bucket().PutObject(pubData.DestLocation, bufio.NewReader(filed))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	//4。更新文件的存储路径到文件表
	sus := db.UpdateFileLocation(pubData.FileHash, pubData.DestLocation)
	if !sus {
		return false
	}
	return true
}

func main() {
	log.Println("开始监听转移任务队列")
	mq.StartConsume(config.TransOSSQueueName, "transfer_oss", ProcessTransfer)
}
