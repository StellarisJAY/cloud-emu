package util

import (
	"github.com/StellrisJAY/cloud-emu/platform/internal/conf"
	"github.com/bwmarrin/snowflake"
)

func NewSnowflakeGenerator(c *conf.Server) *snowflake.Node {
	node, err := snowflake.NewNode(c.NodeId)
	if err != nil {
		panic(err)
	}
	return node
}
