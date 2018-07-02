package common

//LeakyBuf struct
type LeakyBuf struct {
	bufSize  int
	freeList chan []byte
}

const leakyBufSize = 4096
const maxNBuf = 2048

var leakyBuf = NewLeakyBuf(leakyBufSize, maxNBuf)

//NewLeakyBuf create buffer
func NewLeakyBuf(bufSize, n int) *LeakyBuf {
	return &LeakyBuf{
		bufSize:  bufSize,
		freeList: make(chan []byte, n),
	}
}

//Get get leaky buffer data
func (lb *LeakyBuf) Get() (b []byte) {
	select {
	case b = <-lb.freeList:
	default:
		b = make([]byte, lb.bufSize)
	}
	return
}

//Put put data to leaky buffer
func (lb *LeakyBuf) Put(b []byte) {
	if len(b) != lb.bufSize {
		panic("invalid buffer size that's put into leaky buffer")
	}
	select {
	case lb.freeList <- b:
	default:
	}
	return
}
