package main

import (
	"actor-playground/core"
	"actor-playground/game"
	"actor-playground/msg"
	"actor-playground/persis"
	"actor-playground/util"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/persistence"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
)

func TestActorSuite(t *testing.T) {
	suite.Run(t, &TestActor{})
}

type TestActor struct {
	suite.Suite
	sys         *actor.ActorSystem
	rootContext *actor.RootContext
}

func (t *TestActor) SetupTest() {
	t.sys = actor.NewActorSystem()
	t.rootContext = t.sys.Root
	sv := actor.NewOneForOneStrategy(10, 1000, func(reason interface{}) actor.Directive {
		log.Println("系统内发生错误: ", reason)
		return actor.StopDirective
	})
	t.rootContext.WithGuardian(sv)
}

func (t *TestActor) TestHeroUpgradeLv() {
	//Register(persis.ProviderInstance.GetState())
	pid := core.SpawnGameObject(t.rootContext, game.NewRoleService, "rolesvc")
	resp, err := t.rootContext.RequestFuture(pid, &msg.RegisterRoleRequest{
		Name:   "test_role_1",
		Sex:    1,
		Avatar: 0,
	}, util.DefaultTimeout).Result()
	util.Must(err)
	fmt.Println(resp)
	time.Sleep(time.Second)
	fmt.Printf("persistence: \n%s", persis.ProviderInstance.String())
}

func Init(state persistence.ProviderState) {
	hero1 := []byte(`
{
	"id": "hero1",
	"parent_id": "role1",
	"no":	   1,
	"lv":	   1,
	"power":	10
}
`)
	state.PersistSnapshot("role1/hero1", 0, &persis.Snapshot{Data: hero1})
	role1 := []byte(`
{
	"id": 		"role1",
	"lv":	   	1,
	"name":		"kiu",
	"exp":		12,
	"power":	10,
	"hero_list": ["hero1"]
}
`)
	state.PersistSnapshot("role1", 0, &persis.Snapshot{Data: role1})
}
