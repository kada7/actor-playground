/*
 * @author: Haoyuan Liu
 * @date: 2021/1/12
 */

package msg

import "actor-playground/core"

// 添加角色经验
type AddRoleExp struct {
	core.MessageBase
	Exp int
}

// 提升英雄等级
type UpgradeHeroLv struct {
	core.MessageBase
	HeroId string
	AddLv  int
}

// 注册角色
type RegisterRoleRequest struct {
	core.MessageBase
	Name   string //昵称
	Sex    uint8  // 性别：1男2女
	Avatar int    // 形象编号
}

type RegisterRoleResponse struct {
	core.MessageBase
	IsSuccess bool
	Err       error
}

type UnlockHeroRequest struct {
	core.MessageBase
	RoleId string
	No     int
}
type UnlockHeroResp struct {
	core.MessageBase
	HeroId    string
	HeroPower int64
}

// 角色名是否存在请求
type RoleNameExistRequest struct {
	core.MessageBase
	Name string
}

// 角色名是否存在响应
type RoleNameExistResponse struct {
	core.MessageBase
	Existed bool
}
