/*
 * @author: Haoyuan Liu
 * @date: 2021/1/12
 */

package persis_provider

import (
	"actor-playground/core/persistence"
)

//var ProviderInstance = newInMem()
var ProviderInstance = newBolt()

func newInMem() *Provider {
	return NewProvider(persistence.NewInMemoryProvider(3))
}

type Provider struct {
	state persistence.ProviderState
}

func NewProvider(ps persistence.ProviderState) *Provider {
	return &Provider{
		state: ps,
	}
}

func (p *Provider) GetState() persistence.ProviderState {
	return p.state
}
