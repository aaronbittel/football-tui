package component

import (
	"bytes"
	"os"
	"sync"
)

type Buf struct {
	bu    bytes.Buffer
	mutex sync.Mutex
}

func NewBuf() Buf {
	var b bytes.Buffer
	return Buf{
		bu: b,
	}
}

func (b *Buf) Write(s string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.bu.WriteString(s)
}

func (b *Buf) Flush() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.bu.WriteTo(os.Stdout)
}
