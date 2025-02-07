# XcXcStorage

## 介绍
分布式存储节点

## Prerequisites

- **Golang** 1.23.3 or later
- **Etcd** v3.4.0 or later
- **gRPC-go** v1.38.0 or later
- **protobuf** v1.26.0 or later

## Installation

1.  git clone https://gitee.com/Xccccee/xc-xc-storage.git

2.  go mod tidy

3.  go build -o server .

4. ./server

## 使用说明

1.Configure the etcd service 配置etcd

2.RUN ./server 启动server

3.输入etcd地址

4.输入存储目录（相对目录）
