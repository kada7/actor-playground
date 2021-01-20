/*
 * @author: Haoyuan Liu
 * @date: 2021/1/18
 */

package util

import (
	"go.uber.org/atomic"
	"strconv"
)

var uuidNumber = atomic.NewInt64(0)

func NewUUID() string {
	n := uuidNumber.Add(1)
	return strconv.FormatInt(n, 10)
}
