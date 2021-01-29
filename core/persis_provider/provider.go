/*
 * @author: Haoyuan Liu
 * @date: 2021/1/12
 */

package persis_provider

import (
	"actor-playground/core/persistence"
)

// 全局Provider实例
var ProviderInstance = newInMem()

//var ProviderInstance = newBolt("my.db")

func newInMem() *Provider {
	return NewProvider(newInmemProviderState())
}

func newBolt(filename string) *Provider {
	return NewProvider(newBoltProvider(filename))
}

// 自定义的ProviderState
type MyProviderState interface {
	persistence.ProviderState
	// 所有的Actor名字
	ActorNameList() []string
}

type Provider struct {
	state MyProviderState
}

func NewProvider(state MyProviderState) *Provider {
	return &Provider{
		state: state,
	}
}

func (p *Provider) GetState() persistence.ProviderState {
	return p.state
}

func (p *Provider) ActorNameList() []string {
	return p.state.ActorNameList()
}
