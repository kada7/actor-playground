package actor

import (
	"actor-playground/core"
	"actor-playground/game/msg"
	"actor-playground/util"
	"errors"
	"github.com/AsynkronIT/protoactor-go/actor"
)

type RoleCenterState struct {
	// 已经注册的所有角色的Id
	RoleList []string `json:"registered_role"`
}

func (r RoleCenterState) TableName() string {
	return "global"
}

func (r RoleCenterState) PK() string {
	return "role_center"
}

// 角色中心，用于管理所有角色的数据
type RoleCenter struct {
	*core.GameObject
	*RoleCenterState
	rolePidSet  *actor.PIDSet
	roleNameSet *actor.PID
}

func NewRoleCenter() actor.Actor {
	s := &RoleCenter{
		RoleCenterState: &RoleCenterState{RoleList: []string{}},
		rolePidSet:      actor.NewPIDSet(),
	}
	s.GameObject = core.NewGameObject(s.RoleCenterState)
	return s
}

func (s *RoleCenter) Receive(c actor.Context) {
	switch m := c.Message().(type) {
	case *actor.Started:
		s.roleNameSet = core.SpawnGameObject(c, NewRoleNameSet, "role_name_set")
	case *msg.RegisterRoleRequest:
		existed := s.CheckNameExisted(c, m.Name)
		if existed {
			c.Send(c.Sender(), &msg.RegisterRoleResponse{Err: errors.New("role name existed")})
		}
		s.RegisterRole(c)
	}
}

// 校验角色的名字是否存在
func (s *RoleCenter) CheckNameExisted(c actor.Context, name string) bool {
	res, err := c.RequestFuture(s.roleNameSet, &msg.RoleNameExistRequest{Name: name}, util.DefaultTimeout).Result()
	util.Must(err)
	return res.(*msg.RoleNameExistResponse).Existed
}

// 注册角色
func (s *RoleCenter) RegisterRole(c actor.Context) {
	roleId := "role_" + util.NewUUID()
	pid := core.SpawnGameObject(c, NewRole, roleId)
	s.RoleList = append(s.RoleList, roleId)
	s.rolePidSet.Add(pid)
	c.Forward(pid)
}
