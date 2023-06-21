package main

import (
	"github.com/chengjingtao/gomod-version-lint/cmd"
)

func main() {
	cmd.Execute()

	//ctx := context.TODO()
	//log, _ := zap.NewProduction()
	//ctx = pkg.WithLogger(ctx, log.Sugar())
	//
	//res, err := pkg.BranchContains(ctx, "https://github.com/katanomi/builds", "e93ab8d")
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	//fmt.Println(res)

}
