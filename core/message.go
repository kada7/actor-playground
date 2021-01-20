/*
 * @author: Haoyuan Liu
 * @date: 2021/1/20
 */

package core

type Message interface {
	private()
}

type MessageBase struct{}

func (m MessageBase) private() {}
