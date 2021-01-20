package core

import (
	"actor-playground/core/inmem_persis"
	"actor-playground/core/persistence"
	"actor-playground/util"
	"actor-playground/util/ctxlog"
	"encoding/json"
	"github.com/AsynkronIT/protoactor-go/actor"
)

// GameObjecter 提供以下功能
// - 消息、快照的持久化与重放
//
// 所有游戏的Actor都要实现此接口
type GameObjecter interface {
	// 立刻持久化状态
	PersistState()
	// 立刻持久化消息
	PersistMsg(msg interface{})

	recoveryState(s *Snapshot)
	statePtr() interface{}
}

// GameObject 提供以下功能
// - 消息、快照的持久化与重放
//
// 所有游戏的Actor都要继承此游戏对象
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

// 立刻持久化状态
func (g *GameObject) PersistState() {
	b, err := json.Marshal(g.sp)
	util.Must(err)
	g.PersistSnapshot(&Snapshot{Data: b})
}

// 立刻持久化消息
func (g *GameObject) PersistMsg(msg interface{}) {
	g.Mixin.PersistReceive(jsonMessage{msg})
}

func (g *GameObject) recoveryState(s *Snapshot) {
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
		WithReceiverMiddleware(AutoPersisMiddleware).
		WithReceiverMiddleware(persistence.Using(inmem_persis.ProviderInstance)).
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

// 自动持久化的中间件
func AutoPersisMiddleware(next actor.ReceiverFunc) actor.ReceiverFunc {
	return func(c actor.ReceiverContext, envelope *actor.MessageEnvelope) {
		g, ok := c.Actor().(GameObjecter)
		if !ok {
			next(c, envelope)
			return
		}

		if m, ok := envelope.Message.(Message); ok {
			g.PersistMsg(m)
		}
		switch m := envelope.Message.(type) {
		case *Snapshot:
			ctxlog.Debugf(c, "恢复快照, data: %s", string(m.Data))
			g.recoveryState(m)
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
