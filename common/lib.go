package common

const COMPILED_SIZE = 32 << uintptr(^uintptr(0)>>63)

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}
