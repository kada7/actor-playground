package game

import (
	"actor-playground/config"
	"actor-playground/core"
	"actor-playground/msg"
	"actor-playground/persis"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/router"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

// 角色状态
type RoleState struct {
	Id    string `json:"id"`
	Lv    int    `json:"lv"`
	Name  string `json:"name"`
	Exp   int    `json:"exp"`
	Power int64  `json:"power"`
	// 角色拥有的英雄ID列表
	HeroList []string `json:"hero_list"`
}

var _ core.GameObjecter = (*S)(nil)

type S struct {
	core.GameObjecter
}

type Role struct {
	*core.GameObject
	state  *RoleState
	heroes *actor.PID
}

var _ core.GameObjecter = (*Role)(nil)

func NewRole() actor.Actor {
	r := &Role{state: &RoleState{}}
	r.GameObject = core.NewGameObject(r.state)
	return r
}

func (r *Role) Receive(c actor.Context) {
	switch m := c.Message().(type) {
	case *msg.UpgradeHeroLv:
		c.Forward(r.heroes)
	case *msg.AddRoleExp:
		r.AddExp(m.Exp)
	case *persis.Snapshot:
		r.Recovery(c)
		logrus.Infof("角色状态已恢复 RoleState: %+v\n", r.state)
	case *msg.RegisterRoleRequest:
		r.Register(c, m)
		r.PersistState()
		c.Send(c.Sender(), &msg.RegisterRoleResponse{IsSuccess: true})
		logrus.Infof("角色已经初始化 RoleState: %+v\n", r.state)
	case *msg.UnlockHeroResp:
		r.state.HeroList = append(r.state.HeroList, m.HeroId)
		r.state.Power += m.HeroPower
	default:
	}
}

func (r *Role) Register(c actor.Context, m *msg.RegisterRoleRequest) {
	*r.state = RoleState{
		Id:       c.Self().Id,
		Lv:       1,
		Name:     m.Name,
		Exp:      0,
		Power:    0,
		HeroList: make([]string, 0),
	}
	r.UnlockInitHero(c)
}

func (r *Role) AddExp(n int) {
	r.state.Exp += n
}

// 解锁角色初始拥有的英雄
func (r *Role) UnlockInitHero(c actor.Context) {
	initHero := config.GConfigCenter.InitHero
	// 收集英雄初始化结果
	wg := &sync.WaitGroup{}
	awaitPid := c.Spawn(actor.PropsFromFunc(func(c actor.Context) {
		switch m := c.Message().(type) {
		case *msg.UnlockHeroResp:
			r.state.HeroList = append(r.state.HeroList, m.HeroId)
			r.state.Power += m.HeroPower
			wg.Done()
		}
	}))
	// 解锁初始英雄
	for _, heroNo := range initHero {
		wg.Add(1)
		heroId := "hero_" + strconv.Itoa(heroNo)
		pid := core.SpawnGameObject(c, NewHero, heroId)
		c.RequestFuture(pid, &msg.UnlockHeroRequest{
			RoleId: r.state.Id,
			No:     heroNo,
		}, time.Second).PipeTo(awaitPid)
	}
	wg.Wait()
}

func (r *Role) Recovery(c actor.Context) {
	pids := make([]*actor.PID, 0)
	for _, heroId := range r.state.HeroList {
		pids = append(pids, core.SpawnGameObject(c, NewHero, heroId))
	}
	r.heroes = c.Spawn(router.NewConsistentHashGroup(pids...))
}
