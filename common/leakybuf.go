package common

//LeakyBuf 缓冲区结构
type LeakyBuf struct {
	bufSize  int
	freeList chan []byte
}

const leakyBufSize = 4096
const maxNBuf = 2048

var leakyBuf = NewLeakyBuf(leakyBufSize, maxNBuf)

//NewLeakyBuf 漏桶
func NewLeakyBuf(bufSize, n int) *LeakyBuf {
	return &LeakyBuf{
		bufSize:  bufSize,
		freeList: make(chan []byte, n),
	}
}

//Get 新建缓冲区或取缓冲区数据
func (lb *LeakyBuf) Get() (b []byte) {
	select {
	case b = <-lb.freeList:
	default:
		b = make([]byte, lb.bufSize)
	}
	return
}

//Put 存放缓冲数据
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
