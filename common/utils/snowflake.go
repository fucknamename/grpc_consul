package utils

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	// 创建一个雪花节点
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

func GetSnowID() int64 {
	return node.Generate().Int64()
}

func GetSnowIDString() string {
	return node.Generate().String()
}
