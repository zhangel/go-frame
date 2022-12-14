package hooks

import (
	"context"
	"fmt"
	"sync"

	"github.com/modern-go/reflect2"
	"google.golang.org/grpc"
)

var (
	reflectOptField             reflect2.StructField
	reflectUnaryIntField        reflect2.StructField
	reflectStreamingIntField    reflect2.StructField
	reflectCancelField          reflect2.StructField
	reflectMuField              reflect2.StructField
	reflectCsMgrField           reflect2.StructField
	reflectTargetField          reflect2.StructField
	reflectCurBalancerNameField reflect2.StructField

	inMemoryTag = "InMemory"
)

func init() {
	cc := &grpc.ClientConn{}
	reflectCcType := reflect2.TypeOf(cc).(reflect2.PtrType).Elem().(reflect2.StructType)
	reflectOptField = reflectCcType.FieldByName("dopts")
	optField := reflectOptField.Get(cc)
	reflectOptType := reflect2.TypeOf(optField).(reflect2.PtrType).Elem().(reflect2.StructType)
	reflectUnaryIntField = reflectOptType.FieldByName("unaryInt")
	reflectStreamingIntField = reflectOptType.FieldByName("streamInt")
	reflectCancelField = reflectCcType.FieldByName("cancel")
	reflectMuField = reflectCcType.FieldByName("mu")
	reflectCsMgrField = reflectCcType.FieldByName("csMgr")
	reflectTargetField = reflectCcType.FieldByName("target")
	reflectCurBalancerNameField = reflectCcType.FieldByName("curBalancerName")

	if reflectOptField == nil || reflectUnaryIntField == nil || reflectStreamingIntField == nil || reflectCancelField == nil || reflectMuField == nil || reflectCsMgrField == nil || reflectTargetField == nil || reflectCurBalancerNameField == nil {
		panic("Incompatible grpc, replace google.golang.org/grpc => google.golang.org/grpc v1.26.0 in go.mod !")
	}
}

func MakeConnection(target string, unaryInt grpc.UnaryClientInterceptor, streamInt grpc.StreamClientInterceptor) (cc *grpc.ClientConn, err error) {
	return MakeConnectionWithCloseFn(target, unaryInt, streamInt, func() {})
}

func MakeConnectionWithCloseFn(target string, unaryInt grpc.UnaryClientInterceptor, streamInt grpc.StreamClientInterceptor, closeFn func()) (cc *grpc.ClientConn, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("internal:MakeConnection failed, e = %+v", e)
		}
	}()

	var cancelFn context.CancelFunc = func() {
		if closeFn != nil {
			closeFn()
		}
	}

	cc = &grpc.ClientConn{}
	optField := reflectOptField.Get(cc)
	reflectUnaryIntField.Set(optField, &unaryInt)
	reflectStreamingIntField.Set(optField, &streamInt)
	reflectCancelField.Set(cc, &cancelFn)
	reflectMuField.Set(cc, &sync.RWMutex{})
	reflectCsMgrField.Set(cc, reflectCsMgrField.Type().New())
	reflectTargetField.Set(cc, &target)
	reflectCurBalancerNameField.Set(cc, &inMemoryTag)

	return cc, nil
}

func SetCloseFunc(cc *grpc.ClientConn, closeFn func()) {
	reflectCancelField.Set(cc, &closeFn)
}

func IsMemoryConnection(cc *grpc.ClientConn) bool {
	if tag, ok := reflectCurBalancerNameField.Get(cc).(*string); !ok || tag == nil {
		return false
	} else {
		return *tag == inMemoryTag
	}
}
