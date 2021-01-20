package main

import (
	"actor-playground/core"
	"actor-playground/core/inmem_persis"
	actor2 "actor-playground/game/actor"
	"actor-playground/game/msg"
	"actor-playground/util"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
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
	pid := core.SpawnGameObject(t.rootContext, actor2.NewRoleCenter, "rolecenter")
	resp, err := t.rootContext.RequestFuture(pid, &msg.RegisterRoleRequest{
		Name:   "test_role_1",
		Sex:    1,
		Avatar: 0,
	}, util.DefaultTimeout).Result()
	util.Must(err)
	fmt.Println(resp)
	time.Sleep(time.Second)
	fmt.Printf("persistence: \n%s", inmem_persis.ProviderInstance.String())
}
