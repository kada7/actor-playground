package game

import (
	"actor-playground/core"
	"actor-playground/msg"
	"actor-playground/persis"
	"github.com/AsynkronIT/protoactor-go/actor"
	"log"
)

// 英雄状态
type HeroState struct {
	Id       string `json:"id"`
	ParentId string `json:"parent_id"`
	No       int    `json:"no"`
	Lv       int    `json:"lv"`
	Power    int64  `json:"power"`
	// 名媛加成数值
	BeautyBonusValue int64 `json:"beauty_bonus_value"`
	// 名媛加成倍率
	BeautyBonusRate int64 `json:"beauty_bonus_rate"`
}

type Hero struct {
	*core.GameObject
	*HeroState
	rolePid *actor.PID
}

func NewHero() actor.Actor {
	h := &Hero{HeroState: &HeroState{}}
	h.GameObject = core.NewGameObject(h.HeroState)
	return h
}

func (h *Hero) Receive(c actor.Context) {
	switch m := c.Message().(type) {
	case *persis.Snapshot:
		h.Recovery(c)
	case *msg.UnlockHeroRequest:
		h.Unlock(c, m)
		h.PersistState()
	case *msg.UpgradeHeroLv:
		if m.HeroId != h.HeroState.Id {
			break
		}
	}
}

// 初始化英雄
func (h *Hero) Unlock(c actor.Context, m *msg.UnlockHeroRequest) {
	*h.HeroState = HeroState{
		Id:       c.Self().Id,
		ParentId: m.RoleId,
		No:       m.No,
		Lv:       1,
	}
	h.Power = h.CalcHeroPower()
	h.rolePid = c.Parent()
	c.Send(c.Sender(), &msg.UnlockHeroResp{HeroId: h.Id, HeroPower: h.Power})
}

// 提升英雄等级
func (h *Hero) Upgrade(lvNum int) {
	h.Lv += lvNum
	h.Power = h.CalcHeroPower()
	log.Println("英雄Power已变更: ", h.Power)
}

// 重新计算英雄Power
func (h *Hero) CalcHeroPower() int64 {
	power := (int64(h.Lv)*10 + h.BeautyBonusValue) * (1 + h.BeautyBonusRate)
	return power
}

func (h *Hero) Recovery(c actor.Context) {
	h.rolePid = c.Parent()
}
