package component

import (
	"bytes"
	"os"
	"sync"
	"time"
)

type Buf struct {
	bu      bytes.Buffer
	instrCh <-chan string
	mutex   sync.Mutex
}

// TODO: Maybe make this into a singleton?
func NewBuf(instrCh <-chan string) Buf {
	return Buf{
		bu:      bytes.Buffer{},
		instrCh: instrCh,
	}
}

func (b *Buf) ReadLoop() {
	for instr := range b.instrCh {
		b.Write(instr)
	}
}

func (b *Buf) FlushLoop() {
	go func() {
		ticker := time.NewTicker(time.Millisecond * 30)
		for {
			b.flush()
			<-ticker.C
		}
	}()
}

func (b *Buf) Write(s string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.bu.WriteString(s)
}

func (b *Buf) flush() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.bu.WriteTo(os.Stdout)
}
