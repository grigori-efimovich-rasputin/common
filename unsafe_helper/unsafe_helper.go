package unsafeHelper

import (
	"reflect"
	"unsafe"
)

func StringToByte(s string) []byte {
	sh := *(*reflect.SliceHeader)(unsafe.Pointer(&s))
	sh.Cap = len(s)
	return *(*[]byte)(unsafe.Pointer(&sh))
}

func ByteToString(s []byte) string {
	return *(*string)(unsafe.Pointer(&s))
}
