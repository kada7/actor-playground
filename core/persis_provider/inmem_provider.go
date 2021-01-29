/*
 * @author: Haoyuan Liu
 * @date: 2021/1/29
 */

package persis_provider

import (
	"actor-playground/core/persistence"
	"github.com/golang/protobuf/proto"
	"sync"
)

type inmemProviderState struct {
	persistence.ProviderState
	actorName sync.Map
}

func newInmemProviderState() *inmemProviderState {
	return &inmemProviderState{
		ProviderState: persistence.NewInMemoryProvider(3),
		actorName:     sync.Map{},
	}
}

func (s *inmemProviderState) ActorNameList() []string {
	m := make([]string, 0)
	s.actorName.Range(func(key, value interface{}) bool {
		m = append(m, key.(string))
		return true
	})
	return m
}

func (s *inmemProviderState) PersistEvent(actorName string, eventIndex int, event proto.Message) {
	s.actorName.Store(actorName, true)
	s.ProviderState.PersistEvent(actorName, eventIndex, event)
}

func (s *inmemProviderState) PersistSnapshot(actorName string, snapshotIndex int, snapshot proto.Message) {
	s.actorName.Store(actorName, true)
	s.ProviderState.PersistSnapshot(actorName, snapshotIndex, snapshot)
}
