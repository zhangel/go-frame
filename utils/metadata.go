package utils

import (
	"google.golang.org/grpc/metadata"
)

func CopyMetaData(md metadata.MD) metadata.MD {
	newMD := make(metadata.MD, len(md))
	for k, v := range md {
		newMD[k] = v
	}
	return newMD
}
