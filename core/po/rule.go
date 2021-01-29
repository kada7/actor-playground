/*
 * @author: Haoyuan Liu
 * @date: 2021/1/21
 */

package po

import (
	"reflect"
	"strings"
)

// PO Persistence Object 可持久化对象
type PO interface {
	TableName() string // 表名
	PK() string        // 主键
}

// 将PO转为ActorName
func POToActorName(po PO) string {
	return po.TableName() + ":" + po.PK()
}

// POInfo PO信息
type POInfo struct {
	TableName string
	PK        string
}

func ActorNameToPO(actorName string) POInfo {
	result := strings.Split(actorName, ":")
	if len(result) != 2 {
		panic("wrong actor name")
	}
	return POInfo{
		TableName: result[0],
		PK:        result[1],
	}
}

// 解析结构体，返回含有 `model:"pk"` tag的结构体字段,
// 若不存在该tag，则返回false
func ParsePKField(s interface{}) (string, bool) {
	rt := reflect.TypeOf(s)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		panic("s must be struct or struct ptr")
	}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag, ok := field.Tag.Lookup("model")
		if !ok {
			continue
		}
		if tag == "pk" {
			return field.Name, true
		}
	}
	return "", false
}
