package command

// type OnCommand struct {
// 	Lock sync.RWMutex
// 	Data map[int]bool
// }

// type dataResponses struct {
// 	mutex sync.RWMutex
// 	data  map[int][]byte
// }

// // get take data from buffer by key.
// func (d *dataResponses) get(k int) ([]byte, bool) {
// 	d.mutex.RLock()
// 	defer d.mutex.RUnlock()

// 	data, ok := d.data[k]
// 	return data, ok
// }

// // set put data to buffer by key and value.
// func (d *dataResponses) set(k int, v []byte) {
// 	d.mutex.Lock()
// 	defer d.mutex.Unlock()
// 	d.data[k] = v
// }

// // reset delete data from buffer by key.
// func (d *dataResponses) reset(k int) {
// 	d.mutex.Lock()
// 	defer d.mutex.Unlock()
// 	delete(d.data, k)
// }

// var responses = &dataResponses{
// 	mutex: sync.RWMutex{},
// 	data:  map[int][]byte{},
// }
