/*
 * @author: Haoyuan Liu
 * @date: 2021/1/12
 */

package msg

type Message interface {
	private()
}

type MessageBase struct{}

func (m MessageBase) private() {}

// 添加角色经验
type AddRoleExp struct {
	MessageBase
	Exp int
}

// 提升英雄等级
type UpgradeHeroLv struct {
	MessageBase
	HeroId string
	AddLv  int
}

// 注册角色
type RegisterRoleRequest struct {
	MessageBase
	Name   string //昵称
	Sex    uint8  // 性别：1男2女
	Avatar int    // 形象编号
}

type RegisterRoleResponse struct {
	MessageBase
	IsSuccess bool
	Err       error
}

type UnlockHeroRequest struct {
	MessageBase
	RoleId string
	No     int
}
type UnlockHeroResp struct {
	MessageBase
	HeroId    string
	HeroPower int64
}

// 角色名是否存在请求
type RoleNameExistRequest struct {
	MessageBase
	Name string
}

// 角色名是否存在响应
type RoleNameExistResponse struct {
	MessageBase
	Existed bool
}
