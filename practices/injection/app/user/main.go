package main

import (
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"

	"main/app/user/internal/cmd"
)

func main() {
	var ctx = gctx.GetInitCtx()
	c, err := gcmd.NewFromObject(cmd.Main{})
	if err != nil {
		panic(err)
	}
	c.Run(ctx)
}
