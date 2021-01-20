/*
 * @author: Haoyuan Liu
 * @date: 2021/1/12
 */

package inmem_persis

import (
	"actor-playground/core/persistence"
	"actor-playground/util"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"sync"
)

var ProviderInstance = NewProvider(3)

type Snapshot struct {
	Data []byte
}

func (p *Snapshot) Reset()         {}
func (p *Snapshot) ProtoMessage()  {}
func (p *Snapshot) String() string { return string(p.Data) }

type Provider struct {
	state *providerState
}

func NewProvider(snapshotInterval int) *Provider {
	return &Provider{
		state: newProviderState(persistence.NewInMemoryProvider(snapshotInterval)),
	}
}

func (p *Provider) GetState() persistence.ProviderState {
	return p.state
}

func (p *Provider) String() string {
	m := map[string]interface{}{}
	m["snapshot"] = p.getSnapshotData()
	m["message"] = p.getMessage()
	b, err := json.MarshalIndent(m, "", "  ")
	util.Must(err)
	return string(b)
}

func (p *Provider) getMessage() map[string][]interface{} {
	m := map[string][]interface{}{}
	//mutex := &sync.Mutex{}
	//wg := &sync.WaitGroup{}
	for actName := range p.state.ActorNameSet() {
		//wg.Add(1)
		p.state.GetEvents(actName, 0, 0, func(e interface{}) {
			//mutex.Lock()
			//defer mutex.Unlock()
			l, ok := m[actName]
			if !ok {
				l = make([]interface{}, 0)
			}
			l = append(l, e)
			//wg.Done()
			fmt.Printf("[%s]event: %+v\n", actName, e)
		})
	}
	//wg.Wait()
	return m
}

func (p *Provider) getSnapshotData() map[string][]interface{} {
	m := map[string][]interface{}{}
	for actName := range p.state.ActorNameSet() {
		s, idx, ok := p.state.GetSnapshot(actName)
		if !ok {
			continue
		}
		l, ok := m[actName]
		if !ok {
			l = make([]interface{}, 0)
		}
		l = append(l, map[string]interface{}{
			"snapshot": s.(*Snapshot).String(),
			"idx":      idx,
		})
		m[actName] = l
	}
	return m
}

type providerState struct {
	persistence.ProviderState
	actorName sync.Map
}

func newProviderState(s persistence.ProviderState) *providerState {
	return &providerState{
		ProviderState: s,
		actorName:     sync.Map{},
	}
}

func (s *providerState) ActorNameSet() map[string]bool {
	m := map[string]bool{}
	s.actorName.Range(func(key, value interface{}) bool {
		m[key.(string)] = true
		return true
	})
	return m
}

func (s *providerState) PersistEvent(actorName string, eventIndex int, event proto.Message) {
	s.actorName.Store(actorName, true)
	s.ProviderState.PersistEvent(actorName, eventIndex, event)
}

func (s *providerState) PersistSnapshot(actorName string, snapshotIndex int, snapshot proto.Message) {
	s.actorName.Store(actorName, true)
	s.ProviderState.PersistSnapshot(actorName, snapshotIndex, snapshot)
}
