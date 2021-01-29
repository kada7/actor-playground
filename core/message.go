/*
 * @author: Haoyuan Liu
 * @date: 2021/1/20
 */

package core

// 所有业务消息需要实现本接口
type Message interface {
	private()
}

// 所有消息的基类
type MessageBase struct{}

func (m MessageBase) private() {}

//type jsonMessage struct {
//	m interface{}
//}
//
//func (p jsonMessage) Reset() {}
//
//func (p jsonMessage) String() string {
//	b, err := json.Marshal(p.m)
//	util.Must(err)
//	return string(b)
//}
//
//func (p jsonMessage) ProtoMessage() {}
