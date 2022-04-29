package cache

type ByteView struct {
	B []byte
}

func (v ByteView) Len() int {
	return len(v.B)
}

func (v ByteView) ByteSlice() []byte {
	return v.CloneBytes(v.B)
}

func (v ByteView) String() string {
	return string(v.B)
}

func (v ByteView) CloneBytes(b []byte) []byte {
	ret := make([]byte, len(b))
	copy(ret, b)
	return ret
}
