package transform

func I32m2Im(in map[int32]int32) map[int]int {
	out := make(map[int]int, len(in))
	PutI32m2Im(in, out)

	return out
}

func I3264m2I64m(in map[int32]int64) map[int]int64 {
	out := make(map[int]int64, len(in))
	PutI3264m2I64m(in, out)

	return out
}

func I64m2I3264m(in map[int]int64) map[int32]int64 {
	out := make(map[int32]int64, len(in))
	PutI64m2I3264m(in, out)

	return out
}

func PutI32m2Im(in map[int32]int32, out map[int]int) {
	for k, v := range in {
		out[int(k)] = int(v)
	}
}

func PutI3264m2I64m(in map[int32]int64, out map[int]int64) {
	for k, v := range in {
		out[int(k)] = int64(v)
	}
}

func PutI64m2I3264m(in map[int]int64, out map[int32]int64) {
	for k, v := range in {
		out[int32(k)] = int64(v)
	}
}
