/*
 * @author: Haoyuan Liu
 * @date: 2021/1/20
 */

package _file_persis

import (
	"actor-playground/core/po"
	"github.com/golang/protobuf/proto"
	"sync"
)

type entry struct {
	eventIndex int // the event index right after snapshot
	snapshot   proto.Message
	events     []proto.Message
}

type snapshotJson map[string][]interface{}

type snapshotEntry struct{
	index int
	snapshot []byte
}
/*
snapshot example
{
	"hero_1": [
		{"id": "hero_1", "lv": 0}
		{"id": "hero_1", "lv": 1}
	]
	"hero_2": [
		{"id": "hero_1", "lv": 0}
		{"id": "hero_1", "lv": 1}
	]
}

*/

type FileProvider struct {
	snapshotInterval int
	mu               sync.RWMutex
	store            map[string]*entry // actorName -> a persistence entry
}

func NewFileProvider(snapshotInterval int) *FileProvider {
	return &FileProvider{
		snapshotInterval: snapshotInterval,
		store:            make(map[string]*entry),
	}
}

// loadOrInit returns the existing entry for actorName if present.
// Otherwise, it initializes and returns an empty entry.
// The loaded result is true if the entry was loaded, false if initialized.
func (provider *FileProvider) loadOrInit(actorName string) (e *entry, loaded bool) {
	provider.mu.RLock()
	e, ok := provider.store[actorName]
	provider.mu.RUnlock()

	if !ok {
		provider.mu.Lock()
		e = &entry{}
		provider.store[actorName] = e
		provider.mu.Unlock()
	}

	return e, ok
}

func (provider *FileProvider) Restart() {}

func (provider *FileProvider) GetSnapshotInterval() int {
	return provider.snapshotInterval
}

func (provider *FileProvider) GetSnapshot(actorName string) (snapshot interface{}, eventIndex int, ok bool) {
	info := po.ActorNameToPO(actorName)
	entry, loaded := provider.loadOrInit(actorName)
	if !loaded || entry.snapshot == nil {
		return nil, 0, false
	}
	return entry.snapshot, entry.eventIndex, true
}

func (provider *FileProvider) PersistSnapshot(actorName string, eventIndex int, snapshot proto.Message) {
	entry, _ := provider.loadOrInit(actorName)
	entry.eventIndex = eventIndex
	entry.snapshot = snapshot
}

func (provider *FileProvider) DeleteSnapshots(actorName string, inclusiveToIndex int) {

}

func (provider *FileProvider) GetEvents(actorName string, eventIndexStart int, eventIndexEnd int, callback func(e interface{})) {
	entry, _ := provider.loadOrInit(actorName)
	if eventIndexEnd == 0 {
		eventIndexEnd = len(entry.events)
	}
	for _, e := range entry.events[eventIndexStart:eventIndexEnd] {
		callback(e)
	}
}

func (provider *FileProvider) PersistEvent(actorName string, eventIndex int, event proto.Message) {
	entry, _ := provider.loadOrInit(actorName)
	entry.events = append(entry.events, event)
}

func (provider *FileProvider) DeleteEvents(actorName string, inclusiveToIndex int) {

}
