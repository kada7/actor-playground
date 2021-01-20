/*
 * @author: Haoyuan Liu
 * @date: 2021/1/12
 */

package event

// 英雄Power已升级
type HeroPowerUpdated struct {
	HeroId   string
	OldPower int
	NewPower int
}

// 角色名称已经变更
type RoleNameChanged struct {
	OldName string
	NewName string
}
