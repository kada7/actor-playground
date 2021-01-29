# actor游乐场

本代码仓库用于开发采用actor模型的游戏后端demo。

主要采用 [protoactor-go](https://github.com/AsynkronIT/protoactor-go) 这套actor框架

# 前置知识

## 了解actor模型

要了解actor模型，需要先看一看 [Proto.Actor的官方文档](https://proto.actor/docs/) ，
不想看英文的，可以看一下这里粗糙的[中文翻译](https://github.com/kada7/protoactor-go-doc) 。

学有余力的话，强烈推荐阅读 [《Akka实战》](https://book.douban.com/subject/30218333/) ，
学习java的akka框架，该框架与protoactor-go的用法基本相同，本书较系统地展示了akka的各种功能。
并提供了示例。

## 阅读protoactor-go的官方example

[官方example](https://github.com/AsynkronIT/protoactor-go/tree/dev/_examples) ，
建议clone到本地跑跑看。

## 现状

- 暂时实现了一个简单的游戏业务
- protoactor-go内的persistence包不支持中间件，复制了一份persistence包到本项目内，
  并支持了中间件
- 准备尝试实现数据的持久化以及重放（暂时打算使用bolt的provider）
- bolt的provider我fork了一份，之前的bolt-provider没维护了，有些兼容性问题。
- 希望每个游戏业务的actor都通过组合一个GameObject基类来实现一整套游戏的通用功能，如持久化、重放等。
  游戏业务开发人员不需要自己编写持久化的逻辑