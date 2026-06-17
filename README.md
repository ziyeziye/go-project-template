# GinTemplate

基于 Gin 框架的 Go 后端项目模板，内置定时任务调度、多级缓存、统一错误处理和日志系统。

## 技术栈

| 组件 | 库 |
|------|----|
| HTTP 框架 | [Gin](https://github.com/gin-gonic/gin) v1.12 |
| CLI | [Cobra](https://github.com/spf13/cobra) |
| 配置管理 | [Viper](https://github.com/spf13/viper) |
| 数据库 ORM | [GORM](https://gorm.io) + MySQL |
| 缓存 (Redis) | [go-redis](https://github.com/redis/go-redis) + [go-redis/cache](https://github.com/go-redis/cache) |
| 缓存 (内存) | [cache2go](https://github.com/muesli/cache2go) |
| 定时任务 | [robfig/cron](https://github.com/robfig/cron/v3) |
| 数值精度 | [shopspring/decimal](https://github.com/shopspring/decimal) |
| 时间处理 | [carbon](https://github.com/dromara/carbon) v2 |

## 快速开始

```bash
# 进入后端目录
cd backend

# 配置环境变量
cp .env.example .env

# 启动服务（API 接口 + 定时任务）
go run main.go

# 调试模式（打印配置 + 运行测试任务）
go run main.go test
```

### 环境变量

```
PORT=8080              # 服务端口（必填）
DB_HOST=localhost      # MySQL 地址
DB_PORT=3306           # MySQL 端口
DB_DATABASE=           # 数据库名
DB_USERNAME=           # 数据库用户名
DB_PASSWORD=           # 数据库密码
RDB_HOST=localhost     # Redis 地址
RDB_PORT=6379          # Redis 端口
RDB_PASSWORD=          # Redis 密码
```

配置通过 Viper 读取，支持 `.env` 文件和系统环境变量，后者优先级更高。

> 环境变量**不区分大小写**，如 `DB_HOST` 也可写作 `db_host`。

## 命令说明

### `go run main.go` — 启动服务

同时启动两部分：

- **HTTP API 服务**：默认监听 `:8080`，自动注册所有路由
- **定时任务**：以 goroutine 方式并行运行 Cron 任务

非 Release 模式下会在启动时打印所有路由列表：

```
[GET] http://127.0.0.1:8080/api/hello
```

### `go run main.go test` — 调试模式

执行一次性的调试任务，适合：

- 验证配置是否正确加载
- 运行测试 job 检查业务逻辑
- 打印当前配置（使用 `spew.Dump`）

不会启动 HTTP 服务。

## 目录结构

```
backend/
├── main.go                # 程序入口
├── cmd/                   # Cobra 命令
│   ├── root.go            # 默认命令：启动 API + 定时任务
│   ├── cron.go            # 定时任务调度
│   ├── test.go            # 调试命令
│   └── job/
│       └── testjob/       # 示例定时任务（每 6 秒执行一次）
├── app/                   # 应用核心层
│   ├── server.go          # Gin 服务器配置
│   ├── api/               # HTTP 控制器
│   │   ├── base.go        # 基础控制器（可嵌入）
│   │   └── hello.go        # Hello 示例控制器
│   ├── middleware/
│   │   └── testUser.go    # 用户身份中间件
│   ├── errno/             # 业务错误码
│   ├── request/           # 请求结构体定义
│   ├── resp/              # 统一响应格式
│   └── routers/           # 路由注册
├── common/                # 公共组件
│   ├── db/                # MySQL 数据库连接（GORM）
│   ├── rdb/               # Redis 连接 + 缓存封装
│   ├── cachex/            # 内存缓存封装
│   ├── logx/              # 日志组件（按日期切割）
│   └── utils/             # 工具函数
├── config/                # 配置结构体定义
├── model/                 # GORM 数据模型
├── service/               # 业务逻辑层
├── runtime/               # 运行时日志（自动生成）
├── go.mod / go.sum
├── .env / .env.example
```

## 添加新功能

### 新增 API 路由

```go
// 1. 在 app/api/ 下创建控制器
type UserController struct{ api.BaseController }

func (c UserController) Profile(ctx *gin.Context) {
    resp.JsonOk(ctx, resp.H{"name": "Alice"})
}

// 2. 在 app/routers/router.go 中注册路由
func ApiV0(r *gin.Engine) *gin.Engine {
    v0 := r.Group("/api").Use(middleware.TestUser)
    {
        v0.GET("/hello", api.HelloController{}.Hello)
        v0.GET("/user/profile", api.UserController{}.Profile) // 新增
    }
    return r
}
```

### 新增定时任务

```go
// 1. 在 cmd/job/ 下创建新的 job 包
type MyJob struct{}

func NewMyJob() *MyJob { return &MyJob{} }

func (j *MyJob) Run() {
    log.Println("my job running...")
}

// 2. 在 cmd/cron.go 中注册
c.AddJob("@every 30s", cron.NewChain(
    cron.Recover(cron.DefaultLogger),
    cron.SkipIfStillRunning(cron.DefaultLogger),
).Then(MyJob.NewMyJob()))
```

### 新增数据模型

```go
// 1. 在 model/ 下定义结构体
type Account struct {
    ID        uint            `gorm:"primarykey"`
    Name      string          `gorm:"size:255;not null"`
    Status    int8            `gorm:"default:0"`
    CreatedAt carbon.DateTime `gorm:"index"`
    UpdatedAt carbon.DateTime `gorm:"index"`
}

// 2. 在 cmd/root.go 的 afterInit() 中启用迁移
err := db.Engine().AutoMigrate(
    model.Account{},
)
```

### 新增业务逻辑

```go
// 1. 在 service/ 下创建 Service
var AccountSvr = AccountService{}

type AccountService struct{}

func (s *AccountService) GetByID(id uint) (*model.Account, error) {
    // 业务逻辑
}

// 2. 在控制器中调用
func (c AccountController) Get(ctx *gin.Context) {
    account, err := service.AccountSvr.GetByID(1)
    if err != nil {
        resp.JsonErr(ctx, err)
        return
    }
    resp.JsonOk(ctx, account)
}
```

## 响应格式

所有 API 返回统一的 JSON 结构：

```json
// 成功
{"code": 0, "msg": "success", "data": {...}}

// 业务错误（通过 errno 系统）
{"code": 10500, "msg": "Server internal error", "data": null}

// 自定义错误
{"code": -1, "msg": "自定义错误信息", "data": null}
```

## 错误处理

`errno.Error` 为值类型，通过 `errno.ErrOf()` 统一转换：

```go
// 定义业务错误码
var (
    ServerError   = errno.NewError(10500, "Server internal error")
    Unauthorized  = errno.NewError(100402, "Unauthorized")
)

// 在 service 中返回（值类型，可直接 return）
func (s *AccountService) Find() error {
    return errno.ServerError
}

// 在 controller 中统一处理
if err := service.AccountSvr.Find(); err != nil {
    resp.JsonErr(ctx, err)    // 内部通过 ErrOf 自动识别 errno.Error
    return
}
```

`resp.JsonErr()` 和 `resp.JsonErrWithData()` 接受 `error` 接口，内部通过 `errno.ErrOf()` 自动识别 `errno.Error` 类型并提取错误码，其他 error 则包装为通用错误（code=10403，msg 为 error 的文本）。

## 中间件

| 中间件 | 作用 |
|--------|------|
| `CORS` | 允许所有跨域请求 |
| `CustomRecovery` | 捕获 panic，返回 500 JSON |
| `TestUser` | 从 Header/Query/Form 提取 `user_name`，方便调试 |

## 缓存使用

项目提供两级缓存抽象：

### Redis 缓存（持久化、分布式）

```go
import "go-project-template/common/rdb"

// 泛型缓存：GetOrSet 模式，自动序列化
users, err := rdb.SGet(func() ([]User, error) {
    return fetchUsersFromDB()   // 缓存未命中时回调
}, time.Hour, "cache:users")

// 直接操作
rdb.Cache().Set("key", value, time.Minute)
rdb.Cache().Get("key", &dest)
rdb.Cache().Del("key")
```

### 内存缓存（本地、低延迟）

```go
import "go-project-template/common/cachex"

val, err := cachex.GetOrAdd(func() (string, error) {
    return computeExpensiveValue()
}, time.Minute, "my:key")
```

### Redis 分布式锁

```go
// 简单用法（自动释放）
rdb.RunWithLock("job:sync", 10*time.Second, func() {
    // 确保此函数全局只有一个实例在执行
})

// 进阶用法（手动控制）
release, err := rdb.AcquireLock(ctx, "my:lock", 5*time.Second)
if err != nil {
    return err
}
defer release()
```

## License

MIT
