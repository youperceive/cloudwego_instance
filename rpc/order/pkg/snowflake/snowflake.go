package snowflake

import (
	"log"
	"time"

	"github.com/bwmarrin/snowflake"
)

var SfNode *snowflake.Node

func init() {
	// 第一步：先设置Epoch（必须在NewNode之前！）
	snowflake.Epoch = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

	// 第二步：再创建节点（此时节点会使用上面的Epoch）
	nodeID := int64(778)
	var err error
	SfNode, err = snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf("雪花算法初始化失败：%v", err)
	}
}
