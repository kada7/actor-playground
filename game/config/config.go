package config

// 配置中心示例，偷懒大法
var GConfigCenter = ConfigCenter{
	Hero: []Hero{
		{No: 1, Name: "佛罗多"},
		{No: 2, Name: "山姆"},
		{No: 3, Name: "梅里"},
		{No: 4, Name: "皮聘"},
		{No: 5, Name: "甘道夫"},
		{No: 6, Name: "阿拉贡"},
		{No: 7, Name: "莱格拉斯"},
		{No: 8, Name: "金雳"},
		{No: 9, Name: "波罗莫"},
	},
	InitHero: []int{1, 2, 3, 4},
	RoleLevel: []RoleLevel{
		{Lv: 1, UpgradeNeedExp: 100},
		{Lv: 2, UpgradeNeedExp: 200},
		{Lv: 3, UpgradeNeedExp: 300},
	},
}

// 配置中心
type ConfigCenter struct {
	Hero      []Hero
	InitHero  []int // 初试英雄编号
	RoleLevel []RoleLevel
}

// 角色等级配置
type RoleLevel struct {
	Lv             int // 当前等级
	UpgradeNeedExp int // 当前等级升级需要的经验
}

// 英雄配置
type Hero struct {
	No   int // 英雄编号
	Name string
}

// 英雄等级配置
type HeroLevel struct {
	Lv             int // 当前等级
	UpgradeNeedExp int // 当前等级升级需要的经验
}

// 名媛配置
type Beauty struct {
	No   int // 名媛编号
	Name string
}
