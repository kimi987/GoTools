
# 配置文件/生成的Proto/proto中的字段名字/单例

	_ struct{} `file:"fishing/show.txt"`                // 必须指定，读取该配置的文件名字
	_ struct{} `singleton:"true"`                       // 可选，如果配置了改行，意味着这只有个对象，不是数组什么的
	_ struct{} `proto:"shared_proto.FishingShowProto"`  // 可选，该配置序列化的Proto对象
	_ struct{} `protoconfig:"FishingShow"`              // 可选，该配置序列化的Proto对象在Config中的字段名字

特殊情况：

    如果出现gogen生成后，有包没有引用进来导致编译通不过，可以通过把上面的struct{}改成你要导入的那个包中的任何一个对象就可以了，如：

    type EquipCombineData struct {
        _ struct{}             `file:"combine/equip_combine.txt"`
        _ struct{}             `proto:"shared_proto.EquipCombineDataProto"`
        _ *goods.EquipmentData `protoconfig:"equip_combine"`
	}

每个需要生成的配置文件前面要加上 //gogen:config

```go
    // 钓鱼数据
    //gogen:config
    type FishData struct {
        _ struct{} `file:"fishing/fish.txt"`

        Id uint64 `validator:"int>0"` // 钓鱼id

        Prize      *resdata.Prize // 奖励
        prizeBytes []byte

        IsShow bool `protofield:"-"` // 是否钓鱼展示

        ShowContent string // 钓鱼展示内容

        Weight uint64 `validator:"int>0",protofield:"-"` // 权重
    }
```

# 配置的参数

    validator, type, default, head, protofield
    这几个参数的配置格式是用空格分隔，
    如：`validator:"int" type:"enum" default:"3" head:"uid" protofield:"Id,config.U64ToI32(%s)"`

## validator

    validator所有的参数要放在一起，分号分隔，如: `validator:"int,duplicate,notnil"`

---

### 数据验证

|验证         |    说明                  |
|:----:        |:----:                    |
|int         | 整数（不支持小数）             |
|uint        | 正整数（不支持小数）>=0         |
|int>=0      | 正整数（不支持小数）>=0         |
|>=0         | 正整数（不支持小数）>=0         |
|int>0       | 正整数（不支持小数）>0          |
|>0          | 正整数（不支持小数）>0          |
|string      | 任意字符串                 |
|string>0    | 非空字符串                 |
|float64     | 数字（支持小数）              |
|float64>=0  | 正数（支持小数）>=0           |
|float64>0   | 正数（支持小数）>0            |
|bool        | 布尔值，0|1|true|false    |

举例：
```go
type Example struct {
    Id uint64 `validator:"int>0"` // Id必须是正整数（不支持小数）>0
}
```

也可以自己写正则表达式,
举例：
```go
type Example struct {
    Id uint64 `validator:"^(0|[-]?[1-9][0-9]*)$"` // 此处同上面的int
}
```

---

### 允许数组中的数据重复

配置格式(大小写无关)：

    `validator:"d"` `validator:"dup"` `validator:"duplicate"`

举例：
```go
type Example struct {
    Array1 []uint64 `validator:"d"`         // 表示Array1中的数据可以有重复
    Array2 []uint64 `validator:"dup"`       // 表示Array2中的数据可以有重复
    Array3 []uint64 `validator:"duplicate"` // 表示Array3中的数据可以有重复
}
```

---

### 不允许为空(默认不可以为空)

配置格式(大小写无关)：

    `validator:"notnil"`, `validator:"not nil"`, `validator:"notnull"`, `validator:"not null"`

举例：
```go
type Example struct {
    Obj Object `"notnil"` // 表示Obj不可以为空
}
```

---

### 不允许数组中所有的都为空

配置格式：

    `validator:"notallnil"`, `validator:"not all nil"`, `validator:"notallnull"`, `validator:"not all null"`

举例：
```go
type Example struct {
    Objs []Object `validator:"notallnil"`   // 表示Objs中的数据不可以全部为空
}
```

---

### 数组中的数据不允许部分为空(要么都为空，要么都不为空)

配置格式：

    `validator:"allnilornot"`

举例：
```go
type Example struct {
    Objs []Object `validator:"allnilornot"`   // 表示Objs中的数据要么都为空，要么都不为空
}
```

---

### 大小写匹配(默认大小写不匹配)

配置格式：

    `validator:"case"`

举例：
```go
type Example struct {
    Objs uint64 `validator:"case"`  // 表示参数的名字必须叫做Objs, objs也不行
}
```

---

### tips

配置格式：

    `validator:"tips"`

举例：
```go
type Example struct {
    Objs uint64 `validator:"tips:我是tips"`  // 表示报错时的tips内容是 我是tips
}
```

---

### 配置的列的数量

配置格式：

    `validator:"count"`

举例：
```go
type Example struct {
    Objs []uint64 `validator:"count:4"`  // 表示参数的Objs的长度必须等于4
}
```

---

### 数组求和

配置格式：

    `validator:"sum"`

举例：
```go
type Example struct {
    Objs []uint64 `validator:"sum:233"`  // 表示数组之和必须==233
}
```

---

### whiteList参数

    TODO

---

## type参数:

配置格式：

    `type:"enum"` `type:"sub"`

举例:

enum:

```go
    //gogen:config
    type ResAmount struct {
        Type    shared_proto.ResType `type:"enum"`
    }
```

sub: 下面的GoodsEffect在初始化的时候，不会通过GetId之类的去获取，而是直接使用初始化GoodsData的ObjectParser去初始化

```
    //gogen:config
    type GoodsData struct {
        GoodsEffect *GoodsEffect                 `type:"sub"`
    }
```

---

## head参数:

配置格式：

`head:"-"`

    配置 head:"-"，表示该字段不自动初始化

`head:"equip_id"`

    配置 head:"equip_id"，表示该字段读取表中的equip_id字段来初始化

`head:"-,指定方法生成"`

    如：Id uint64   `head:"-,DonateId(%s.Sequence%c %s.Times)" protofield:"-"`
    表示 Id 是通过 DonateId 方法传进去参数生成的

---

## protofield参数:

    指定在Proto(pb.go, 不是.proto中的字段名字)中的字段名字，以及指定proto字段在赋值的时候的数据组成形式

配置格式:

    EquipId uint // 不填写，如果该字段是大写字母开头的话，那么默认用该字段名字作为Proto的字段，如:proto.EquipId = Data.EquipId

    EquipId uint `protofield:"-"` // 不序列化给Proto

    EquipId uint `protofield:"EquipId"` // 序列化给Proto的时候，Proto的字段为EquipId，如：Proto.EquipId = Data.EquipId

    EquipId uint `protofield:",config.U64ToI32(%s)"` // 序列化给Proto的时候，Proto的字段为EquipId，如：Proto.EquipId = config.U64ToI32(Data.EquipId)

    EquipId uint `protofield:"EquipmentId,config.U64ToI32(%s)"` // 序列化给Proto的时候，Proto的字段为EquipId，如：Proto.EquipmentId = config.U64ToI32(Data.EquipId)

    EquipData []*goods.EquipmentData `protofield:",config.U64a2I32a(goods.GetEquipmentDataKeyArray(%s))"` // 序列化给Proto的时候，Proto的字段为EquipId，如：Proto.EquipId = config.U64a2I32a(goods.GetEquipmentDataKeyArray(data.EquipData))

    EquipData []*goods.EquipmentData `protofield:"EquipmentId,config.U64a2I32a(goods.GetEquipmentDataKeyArray(%s))"` // 序列化给Proto的时候，Proto的字段为EquipId，如：Proto.EquipmentId = config.U64a2I32a(goods.GetEquipmentDataKeyArray(data.EquipData))

---

## default参数:

    指定默认值

配置格式:

    id uint `default:"3"` // 默认值为3

    name string `default:"sss"` // 默认值为"sss"

    name string `default:"null"` // 默认值为空
