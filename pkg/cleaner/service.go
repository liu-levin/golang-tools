package cleaner

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/liu-levin/golang-tools/pkg/logger"
)

var (
	commandTemplate = `find %s -name "*.JPG" -type f -mtime +60  -exec rm {} \;`
)

type Cleaner interface {
	Run()
}

func NewCleaner(filePath string) Cleaner {

	if filePath == "" {
		filePath = "."
	}

	return &cleaner{
		filePath: filePath,
	}
}

type cleaner struct {
	filePath    string
	expiredDays int
}

func (c *cleaner) Run() {
	var (
		err    error
		out    bytes.Buffer
		stderr bytes.Buffer
	)
	logger.Info.Println("#### start clean job")
	ctx := context.Background()
	cmdStr := fmt.Sprintf(commandTemplate, c.filePath)
	logger.Info.Println("exec is ", cmdStr)
	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		logger.Error.Printf("cmd err is %+v\n", err)
	}
}
