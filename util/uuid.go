/*
 * @author: Haoyuan Liu
 * @date: 2021/1/18
 */

package util

import (
	"github.com/google/uuid"
	"go.uber.org/atomic"
)

var uuidNumber = atomic.NewInt64(0)

// new inmem uuid
//func NewUUID() string {
//	n := uuidNumber.Add(1)
//	return strconv.FormatInt(n, 10)
//}

func NewUUID() string {
	return uuid.New().String()
}
