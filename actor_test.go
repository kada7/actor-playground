package main

import (
	"actor-playground/core"
	"actor-playground/core/persis_provider"
	proto2 "actor-playground/core/proto"
	actor2 "actor-playground/game/actor"
	"actor-playground/game/msg"
	"actor-playground/util"
	"encoding/json"
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

func (t *TestActor) TearDownSuite() {
	fmt.Printf("persistence: \n%s", ProviderString(persis_provider.ProviderInstance))
	persis_provider.CloseDB()
	//persis_provider.CleanBolt()
}

func (t *TestActor) TestJustRun() {
	pid := core.SpawnGameObject(t.rootContext, actor2.NewRoleCenter, "rolecenter")
	defer func() {
		t.rootContext.PoisonFuture(pid).Wait()
	}()
	time.Sleep(5 * time.Second)
}

func (t *TestActor) TestHeroUpgradeLv() {
	pid := core.SpawnGameObject(t.rootContext, actor2.NewRoleCenter, "rolecenter")
	defer func() {
		t.rootContext.PoisonFuture(pid).Wait()
	}()
	resp, err := t.rootContext.RequestFuture(pid, &msg.RegisterRoleRequest{
		Name:   "test_role_1",
		Sex:    1,
		Avatar: 0,
	}, util.DefaultTimeout).Result()
	util.Must(err)
	fmt.Println(resp)
	time.Sleep(time.Second)
}

func ProviderString(p *persis_provider.Provider) string {
	m := map[string]interface{}{}
	m["snapshot"] = getSnapshotData(p)
	m["message"] = getMessage(p)
	b, err := json.MarshalIndent(m, "", "  ")
	util.Must(err)
	return string(b)
}

func getMessage(p *persis_provider.Provider) map[string][]interface{} {
	m := map[string][]interface{}{}
	//mutex := &sync.Mutex{}
	//wg := &sync.WaitGroup{}
	keys, err := persis_provider.GetAllKeys()
	util.Must(err)
	for _, actName := range keys {
		//wg.Add(1)
		p.GetState().GetEvents(actName, 0, 0, func(e interface{}) {
			//mutex.Lock()
			//defer mutex.Unlock()
			l, ok := m[actName]
			if !ok {
				l = make([]interface{}, 0)
			}
			l = append(l, e)
			//wg.Done()
			fmt.Printf("actor[%s]event: %+v\n", actName, e)
		})
	}
	//wg.Wait()
	return m
}

func getSnapshotData(p *persis_provider.Provider) map[string][]interface{} {
	m := map[string][]interface{}{}
	keys, err := persis_provider.GetAllKeys()
	util.Must(err)
	for _, actName := range keys {
		s, idx, ok := p.GetState().GetSnapshot(actName)
		if !ok {
			continue
		}
		l, ok := m[actName]
		if !ok {
			l = make([]interface{}, 0)
		}
		l = append(l, map[string]interface{}{
			"snapshot": s.(*proto2.Snapshot).String(),
			"idx":      idx,
		})
		m[actName] = l
	}
	return m
}
