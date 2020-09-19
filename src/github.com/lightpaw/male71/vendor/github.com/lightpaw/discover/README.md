# Zookeeper Recipe

提供节点注册, 并监听路径下所有的节点, 监听新增/更新/删除事件

## 连接到zk

```go
	conn, err := Connect([]string{"localhost:2181", "localhost:2182", "localhost:2183"}, 4*time.Second, connListener)
```

其中`connListener`是监听当前到Zk的Session是否Expire的回调interface, 需实现

```go
    type ConnectionStateCallback interface {
        // 与zk连接expire时调用, 只会调用一次
        OnSessionExpired()
    }
```

连接到zk成功后`Connect`方法才会返回. 

## 注册

```go
	result, err := conn.CreateEphemeralSequential("/w/server/t-", []byte{1, 2, 3})
```

## 监听

```go
	_, eventChan := conn.Discover("/w/server")
```

`eventChan`中会包含这个路径下所有的新增/更新/删除事件. 
