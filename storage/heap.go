package storage

import (
	app "github.com/ciphermarco/BOAST"
	"github.com/ciphermarco/BOAST/log"
)

type eventHeap []app.Event

func (h eventHeap) Len() int      { return len(h) }
func (h eventHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h eventHeap) Less(i, j int) bool {
	return h[i].Time.Before(h[j].Time)
}

func (h *eventHeap) Push(x interface{}) {
	v, ok := x.(app.Event)
	if !ok {
		log.Info("An error occurred and an event could not be pushed to the events heap")
		log.Debug("eventHeap.Push got data of type %T but wanted boast.Event", v)
	} else {
		*h = append(*h, v)
	}
}

func (h *eventHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
