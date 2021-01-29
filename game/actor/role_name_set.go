package actor

import (
	"actor-playground/core"
	"actor-playground/game/event"
	"actor-playground/game/msg"
	"github.com/AsynkronIT/protoactor-go/actor"
)

type RoleNameSetState struct {
	Set map[string]bool
}

func (r RoleNameSetState) TableName() string {
	return "global"
}

func (r RoleNameSetState) PK() string {
	return "role_name_set"
}

// 角色的名字集合
type RoleNameSet struct {
	*core.GameObject
	*RoleNameSetState
}

func NewRoleNameSet() actor.Actor {
	s := &RoleNameSet{RoleNameSetState: &RoleNameSetState{
		Set: map[string]bool{},
	}}
	s.GameObject = core.NewGameObject(s.RoleNameSetState)
	return s
}

func (r *RoleNameSet) Receive(c actor.Context) {
	switch m := c.Message().(type) {
	case *actor.Started:
		r.SubscriberRoleNameChanged(c)
	case *msg.RoleNameExistRequest:
		existed := r.IsExist(m.Name)
		c.Send(c.Sender(), &msg.RoleNameExistResponse{Existed: existed})
	}
}

// 订阅角色名称已变更事件，替换自身集合内的数据
func (r *RoleNameSet) SubscriberRoleNameChanged(c actor.Context) {
	c.ActorSystem().EventStream.Subscribe(func(evt interface{}) {
		e := evt.(*event.RoleNameChanged)
		if e.OldName != "" {
			delete(r.Set, e.OldName)
		}
		r.Set[e.NewName] = true
	}).WithPredicate(func(evt interface{}) bool {
		_, ok := evt.(*event.RoleNameChanged)
		return ok
	})
}

// 名称是否存在
func (r *RoleNameSet) IsExist(name string) bool {
	_, ok := r.Set[name]
	return ok
}
