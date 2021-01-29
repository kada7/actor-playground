/*
 * @author: Haoyuan Liu
 * @date: 2021/1/12
 */

package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"reflect"
	"time"
)

var DefaultTimeout = time.Hour

func CopyIntList(src []int) []int {
	to := make([]int, 0)
	copy(to, src)
	return to
}

func SendMany(c actor.Context, pids []*actor.PID, message interface{}) {
	for i := range pids {
		c.Send(pids[i], message)
	}
}

func SendMany2(c actor.Context, set *actor.PIDSet, message interface{}) {
	set.ForEach(func(i int, pid *actor.PID) {
		c.Send(pid, message)
	})
}

func Must(err error) {
	if err != nil {
		fmt.Println("Must err", err)
		panic(err)
	}
}

func JsonMarshal(v interface{}) string {
	b, err := json.Marshal(v)
	Must(err)
	return string(b)
}

func JsonUnmarshal(s string, v interface{}) {
	b := bytes.NewBufferString(s)
	err := json.Unmarshal(b.Bytes(), v)
	Must(err)
}

// 获取结构体的名字
func GetStructName(s interface{}) string {
	name, err := GetStructName2(reflect.TypeOf(s))
	Must(err)
	return name
}

func GetStructName2(s reflect.Type) (string, error) {
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	if s.Kind() != reflect.Struct {
		return "", errors.New("s must be struct or struct ptr")
	}
	return s.Name(), nil
}
