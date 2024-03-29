package events

import (
	"reflect"
	"sort"
	"sync"
)

type bitmask uint8
type Handler func()

const (
	ExitEvent bitmask = 1 << iota
	StopEvent

	totalEvents int = 2
)

type order int8

const (
	FirstOrder order = iota - 1
	NormalOrder
	LastOrder
)

type handlerDatum struct {
	event     bitmask
	order     order
	handler   Handler
	isRunning bool
	mutex     sync.Mutex
}

type handlerData []*handlerDatum

func (d handlerData) Len() int           { return len(d) }
func (d handlerData) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d handlerData) Less(i, j int) bool { return (d[j].order - d[i].order) > 0 }

var (
	data      = make(handlerData, 0, 1<<10)
	mutex     = sync.RWMutex{}
	isExiting = false
)

func IsExiting() bool {
	defer mutex.RUnlock()
	mutex.RLock()

	return isExiting
}

func getHandlers(e bitmask, localData *handlerData) {
	defer mutex.RUnlock()
	mutex.RLock()

	l := len(data)
	*localData = make(handlerData, 0, l)

	for i := 0; i < l; i++ {
		if datum := data[i]; datum != nil {
			for j := 0; j < totalEvents; j++ {
				if b := e & (1 << bitmask(j)); datum.event&b > 0 {
					*localData = append(*localData, datum)
					break
				}
			}
		}
	}
}

func RunHandlers(e bitmask) {
	var localData handlerData

	getHandlers(e, &localData)

	sort.Sort(localData)

	for _, datum := range localData {
		var run bool

		func() {
			defer datum.mutex.Unlock()
			datum.mutex.Lock()

			if !datum.isRunning {
				run, datum.isRunning = true, true
			}
		}()

		if !run {
			continue
		}

		func() {
			defer func() {
				datum.isRunning = false
				datum.mutex.Unlock()
			}()
			datum.mutex.Lock()

			datum.handler()
		}()
	}
}

func AddHandler(event bitmask, order order, handler Handler) {
	defer mutex.Unlock()
	mutex.Lock()

	// just update event bitmask if handler function is already found
	for _, datum := range data {
		if datum.order == order {
			p1 := reflect.ValueOf(handler).Pointer()
			p2 := reflect.ValueOf(datum.handler).Pointer()

			if p1 == p2 {
				datum.event |= event
				return
			}
		}
	}

	datum := handlerDatum{
		event:   event,
		order:   order,
		handler: handler,
	}

	data = append(data, &datum)
}

func AddNormalHandler(event bitmask, handler Handler) {
	AddHandler(event, NormalOrder, handler)
}

func init() {
	AddNormalHandler(
		ExitEvent,
		func() {
			defer mutex.Unlock()
			mutex.Lock()

			isExiting = true
		},
	)
}
