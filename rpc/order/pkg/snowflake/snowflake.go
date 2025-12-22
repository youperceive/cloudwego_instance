package snowflake

import (
	"log"
	"time"

	"github.com/bwmarrin/snowflake"
)

var SfNode *snowflake.Node

func init() {
	nodeID := int64(1)
	var err error
	SfNode, err = snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf("雪花算法初始化失败：%v", err)
	}

	snowflake.Epoch = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
}
