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
![输入图片说明](%E5%B1%8F%E5%B9%95%E6%88%AA%E5%9B%BE%202025-02-07%20211300.png)

5.添加节点
![输入图片说明](%E5%B1%8F%E5%B9%95%E6%88%AA%E5%9B%BE%202025-02-07%20211327.png)

6.节点发现
![输入图片说明](%E5%B1%8F%E5%B9%95%E6%88%AA%E5%9B%BE%202025-02-07%20211331.png)

7.删除节点
![输入图片说明](%E5%B1%8F%E5%B9%95%E6%88%AA%E5%9B%BE%202025-02-07%20211356.png)
