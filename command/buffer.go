package command

import "sync"

// type OnCommand struct {
// 	Lock sync.RWMutex
// 	Data map[int]bool
// }

type OnResponse struct {
	Mutex sync.RWMutex
	Data  map[int][]byte
}

func (d *OnResponse) Get(k int) ([]byte, bool) {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()

	data, ok := d.Data[k]
	return data, ok
}

func (d *OnResponse) Set(k int, v []byte) {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	d.Data[k] = v
}

func (d *OnResponse) Reset(k int) {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	delete(d.Data, k)
}

var RX = &OnResponse{
	Mutex: sync.RWMutex{},
	Data:  map[int][]byte{},
}
