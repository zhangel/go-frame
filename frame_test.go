package framework

import (
	"context"
	"github.com/zhangel/go-frame.git/log"
	"testing"
)

//go test -test.run TestGoFramework
func TestGoFramework(t *testing.T) {
	defer Init()()
	logger:=log.WithContext(context.Background())
	logger.Fatalf("aaa=%s","ccc")
}
