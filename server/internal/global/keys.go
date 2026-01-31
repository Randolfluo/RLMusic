package g

// 上下文键名,用于存储数据库连接等对象
const (
	CtxDB       = "db"        // 数据库连接对象
	CtxRedis    = "redis"     // Redis连接对象
	CtxUserAuth = "user_auth" // 当前登录用户信息
)
