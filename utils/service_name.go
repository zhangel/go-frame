package utils

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

const targetMeta = "x-grpc-target"

func ServerServiceMethodName(ctx context.Context, fullMethod string) (string, string) {
	target := ""

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		targetMd := md.Get(targetMeta)
		if len(targetMd) != 0 {
			target = targetMd[len(targetMd)-1]
		}
	}

	path := strings.Split(fullMethod, "/")
	if len(path) == 3 {
		if target != "" {
			return target, path[2]
		} else {
			return path[1], path[2]
		}
	}

	return "", ""
}

func ClientServiceMethodName(target, fullMethod string) (string, string) {
	path := strings.Split(fullMethod, "/")
	if len(path) == 3 {
		if target != "" && !strings.Contains(target, ":") {
			return target, path[2]
		} else {
			return path[1], path[2]
		}
	}

	return "", ""
}
