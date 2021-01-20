package core

import (
	"actor-playground/core/persistence"
	"actor-playground/msg"
	"actor-playground/persis"
	"actor-playground/util"
	"actor-playground/util/ctxlog"
	"encoding/json"
	"github.com/AsynkronIT/protoactor-go/actor"
)

type GameObjecter interface {
	PersistState()
	PersistMsg(msg interface{})
	RecoveryState(s *persis.Snapshot)
	statePtr() interface{}
}

type GameObject struct {
	persistence.Mixin
	sp interface{}
}

func NewGameObject(state interface{}) *GameObject {
	return &GameObject{sp: state}
}

func (g *GameObject) statePtr() interface{} {
	return g.sp
}

func (g *GameObject) PersistState() {
	b, err := json.Marshal(g.sp)
	util.Must(err)
	g.PersistSnapshot(&persis.Snapshot{Data: b})
}

func (g *GameObject) PersistMsg(msg interface{}) {
	g.Mixin.PersistReceive(jsonMessage{msg})
}

func (g *GameObject) RecoveryState(s *persis.Snapshot) {
	err := json.Unmarshal(s.Data, g.sp)
	util.Must(err)
}

func SpawnGameObject(c actor.SpawnerContext, p actor.Producer, actorId string) *actor.PID {
	pid, err := c.SpawnNamed(gameObjectProps(p), actorId)
	util.Must(err)
	return pid
}

func gameObjectProps(producer actor.Producer) *actor.Props {
	return actor.PropsFromProducer(producer).
		//WithContextDecorator(AutoPersis).
		WithReceiverMiddleware(AutoPersisMiddleware).
		WithReceiverMiddleware(persistence.Using(persis.ProviderInstance)).
		WithReceiverMiddleware(LogMiddleware)
}

// 记录日志的中间件
func LogMiddleware(next actor.ReceiverFunc) actor.ReceiverFunc {
	return func(c actor.ReceiverContext, envelope *actor.MessageEnvelope) {
		m := envelope.Message
		ctxlog.Debugf(c, "接收到消息[%s] msg: %+v", util.GetStructName(m), m)
		switch envelope.Message.(type) {
		case *actor.Started:
			ctxlog.Debug(c, "游戏对象已启动")
		default:
		}
		next(c, envelope)
	}
}

func AutoPersis(next actor.ContextDecoratorFunc) actor.ContextDecoratorFunc {
	return func(c actor.Context) actor.Context {
		g, ok := c.Actor().(GameObjecter)
		if !ok {
			return next(c)
		}
		c = next(c)
		if m, ok := c.Message().(msg.Message); ok {
			g.PersistMsg(m)
		}
		switch m := c.Message().(type) {
		case *persis.Snapshot:
			ctxlog.Debugf(c, "恢复快照, data: %s", string(m.Data))
			g.RecoveryState(m)
		case *persistence.RequestSnapshot:
			ctxlog.Debugf(c, "请求生成快照")
			g.PersistState()
			ctxlog.Debugf(c, "快照生成完成")
		case *persistence.ReplayComplete:
			ctxlog.Debugf(c, "重放快照完成，当前state: %+v", g.statePtr())
		}
		return c
	}
}

// 自动持久化的中间件
func AutoPersisMiddleware(next actor.ReceiverFunc) actor.ReceiverFunc {
	return func(c actor.ReceiverContext, envelope *actor.MessageEnvelope) {
		g, ok := c.Actor().(GameObjecter)
		if !ok {
			next(c, envelope)
			return
		}

		if m, ok := envelope.Message.(msg.Message); ok {
			g.PersistMsg(m)
		}
		switch m := envelope.Message.(type) {
		case *persis.Snapshot:
			ctxlog.Debugf(c, "恢复快照, data: %s", string(m.Data))
			g.RecoveryState(m)
		case *persistence.RequestSnapshot:
			ctxlog.Debugf(c, "请求生成快照")
			g.PersistState()
			ctxlog.Debugf(c, "快照生成完成")
		case *persistence.ReplayComplete:
			ctxlog.Debugf(c, "重放快照完成，当前state: %+v", g.statePtr())
		}
		next(c, envelope)
	}
}

type jsonMessage struct {
	m interface{}
}

func (p jsonMessage) Reset() {}

func (p jsonMessage) String() string {
	b, err := json.Marshal(p.m)
	util.Must(err)
	return string(b)
}

func (p jsonMessage) ProtoMessage() {}
