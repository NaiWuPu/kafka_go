### docker  etcd 安装
```
rm -rf /tmp/etcd-data.tmp && mkdir -p /tmp/etcd-data.tmp &&   docker run -d   -p 2379:2379   -p 2380:2380   --mount type=bind,source=/tmp/etcd-data.tmp,destination=/etcd-data   --name etcd-gcr-v3.3.13   quay.io/coreos/etcd:v3.3.13   /usr/local/bin/etcd   --name s1   --data-dir /etcd-data   --listen-client-urls http://0.0.0.0:2379   --advertise-client-urls http://0.0.0.0:2379   --listen-peer-urls http://0.0.0.0:2380   --initial-advertise-peer-urls http://0.0.0.0:2380   --initial-cluster s1=http://0.0.0.0:2380   --initial-cluster-token tkn   --initial-cluster-state new
```

#### 本地测试使用的kafka 是 1.19.0
#### centos 上的环境 kafka version: 1.1.0

```MakeFile
version: '3'

services:
  zookeeper:
    image: wurstmeister/zookeeper
    hostname: zookeeper
    ports:
      - "2181:2181"
    container_name: zookeeper

  
  kafka1:
    image: wurstmeister/kafka
    depends_on: [ zookeeper ]
    ports:
      - "9092:9092"
    environment:
      KAFKA_ADVERTISED_HOST_NAME: localhost
      KAFKA_ZOOKEEPER_CONNECT: "zookeeper:2181"
      KAFKA_BROKER_ID: 1
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_CREATE_TOPICS: "test:1:1"
    container_name: kafka
```

***
#### 写入文件 docker-compose.yml 
#### 并且 docker-compose up -d

***
### 目录结构 
```tree
├─conf      // 配置文件夹
├─etcd      // 从etcd 读取数据 watched etcd 变化
├─kafka     // kafka初始化 对外暴露接口 发送数据到kafka
├─log       // 模拟的日志文件夹
├─software  // 包内含有 windows 下的ETCD kafka zookeeper 运行顺序 zookper.bat kafka.bat
├─tailLog   // 日志收集模块 从文件中持续读取数据
├─utils     // 获取当前ip
└─vendor    // 三方包

```