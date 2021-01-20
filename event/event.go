/*
 * @author: Haoyuan Liu
 * @date: 2021/1/12
 */

package event

type HeroPowerUpdated struct {
	HeroId string
}

// 角色名称已经变更
type RoleNameChanged struct {
	OldName string
	NewName string
}
