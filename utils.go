package bip39go

import "C"

//#include <stdlib.h>
import "C"
import "unsafe"

func memset(a []uint8, val uint8) {
	if len(a) == 0 {
		return
	}
	a[0] = val
	for bp := 1; bp < len(a); bp *= 2 {
		copy(a[bp:], a[:bp])
	}
}

func memsetString(a []string, val string) {
	if len(a) == 0 {
		return
	}
	a[0] = val
	for bp := 1; bp < len(a); bp *= 2 {
		copy(a[bp:], a[:bp])
	}
}

// This is very unsafe option, be absolutely sure you have a valid size of array
func fromCStringArray(len C.size_t, in **C.char) []string {
	var out = make([]string, uint64(len))
	var offset = unsafe.Sizeof(uintptr(0))

	// don't try to get out of bounds value
	for n := uint64(0); n < uint64(len); n++ {
		out[n] = C.GoString(*in)
		// move pointer
		in = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(in)) + offset))
	}

	return out
}

func toCArrayByte(in []uint8, out []C.uchar) {
	for idx, val := range in {
		out[idx] = C.uchar(val)
	}
}

func toCArrayBytePtr(in []uint8, out *C.uchar) {
	var offset = unsafe.Sizeof(uint8(0))

	for n := 0; n < len(in); n++ {
		*out = C.uchar(in[n])
		// move pointer
		out = (*C.uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(out)) + offset))
	}
}

func fromCArrayByte(in []C.uchar, out []uint8) {
	for idx, val := range in {
		out[idx] = uint8(val)
	}
}

func fromCPtrArray(in *C.uchar, out []uint8) {
	var offset = unsafe.Sizeof(C.uchar(0))

	// don't try to get out of bounds value
	for n := 0; n < len(out); n++ {
		out[n] = uint8(*in)
		// move pointer
		in = (*C.uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(in)) + offset))
	}
}

func NewCByteArray(len uint64) *C.uchar {
	return (*C.uchar)(C.malloc(C.size_t(len) * C.size_t(unsafe.Sizeof(uintptr(0)))))
}

func (target *C.uchar) Write(data []uint8) {
	toCArrayBytePtr(data, target)
}

func (target *C.uchar) Free() {
	C.free(unsafe.Pointer(target))
}
