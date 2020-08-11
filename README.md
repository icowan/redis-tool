# redis-tool Redis迁移工具

在平时工作有可能会遇到全并Redis、拆分Redis、Redis单点到集群的迁移问题。

这里写了一个简单的迁移工具，支持多平台

## 支持的数据类型

- [x] `string` 字符串
- [x] `hash` 散列列表
- [x] `list` 表列
- [x] `sorted-set` 有序集合
- [x] `all` 所有

## 命令使用支持

- all         迁移所有
- hash        哈希列表迁移
- set         redis string  迁移
- sorted-set  有序集合迁移
- list        列表

## 使用教程

```
数据迁移命令

Usage:
  redis-tool migrate [command]

Examples:

支持命令:
[hash, set, sorted-set, all]

Available Commands:
  all         迁移所有
  hash        哈希列表迁移
  set         redis set 迁移
  sorted-set  有序集合迁移

Flags:
  -h, --help                   help for migrate
      --source-auth string     源密码
      --source-database int    源database
      --source-hosts string    源redis地址, 多个ip用','隔开 (default "127.0.0.1:6379")
      --source-prefix string   源redis前缀
      --source-redis-cluster   源redis是否是集群
      --target-auth string     目标密码
      --target-database int    目标database
      --target-hosts string    目标redis地址, 多个ip用','隔开 (default "127.0.0.1:6379")
      --target-prefix string   目标redis前缀
      --target-redis-cluster   目标redis是否是集群

Use "redis-tool migrate [command] --help" for more information about a command.
```

### 编译

```bash
make build
```

### 执行

```bash
redis-tool migrate
```

#### 迁移hash

> 单点到集群

```bash
redis-tool migrate hash sys:user --source-hosts=127.0.0.1:6379 --source-auth=123456 --source-database=1 --target-redis-cluster=true --target-hosts=127.0.0.1:6379,127.0.0.1:7379 --target-auth=123456
```

> 集群到集群

```bash
redis-tool migrate hash sys:user --source-hosts 127.0.0.1:6379,127.0.0.1:7379  --source-redis-cluster true --source-auth 123456 --target-redis-cluster true --target-hosts 127.0.0.1:6379,127.0.0.1:7379 --target-auth 123456
```
