# IPFS-DHT-爬虫服务（网站）

演示网址

http://185.239.68.188/

### 编译

$ go build server.go

源码是基于gin框架

### 功能介绍

创建一个http服务，监听爬虫提交的数据，保存到数据库。并在网页展示。爬虫源码暂不开源。

### 运行

运行网站

$ ./server

运行缓存更新工具

https://github.com/magnshen/IPFS-DHT-Spider-WebsiteCache   //未来会把这个做到本服务里

运行爬虫

https://github.com/magnshen/IPFS-DHT-Spider

##### 提交接口：

地址：{host}:{port}/api/submit

方法：post

数据：{

​	"hashs":["QmckJaaXTMUbAdKaom4mvntAichDkwkb4AqjhKwLnGLsmm",......], //哈希值数组

​	"nodeId":"QmUJmrjqh23T5uFgJreeYgNP3ThhHFGdeoSRFjVujCcwuZ", //节点id

​	"spiderName":"xxxxx"//爬虫的名字，可以知道哪个爬虫没在工作

​	}

返回：Success

### 数据库

数据库创建脚本 CreateDB.sql

此数据库也是IPFS-DHT-Spide-Website 网站使用的数据库。

config.cnf 配置文件中的数据库链接配置是指当前数据库