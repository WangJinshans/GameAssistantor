package main

import (
	"github.com/casbin/casbin/v2"
	casbin_model "github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"
)

// 参考: https://casbin.org/zh/docs/get-started/

// RBAC多域名(租户)模式
func CasbinRbacWithDomain() {

	engine, err := gorm.Open(sqlite.Open("debug/casbin/gorm_rbac_domain.db"), &gorm.Config{})
	if err != nil {
		log.Info().Msgf("create db error: %v", err)
		return
	}

	// 字符串模式
	// policy := `
	// p, alice, data1, read
	// p, bob, data2, write
	// p, data2_admin, data2, read
	// p, data2_admin, data2, write
	// g, alice, data2_admin
	// `
	// policy := NewStringAdapter(policy)
	policy, err := gormadapter.NewAdapterByDB(engine)
	if err != nil {
		log.Info().Msgf("init auth error: %v", err)
		return
	}

	// enforcer, err := casbin.NewEnforcer("./auth_model.conf", policy)
	m, err := casbin_model.NewModelFromString(`[request_definition]
	r = sub, dom, obj, act
	
	[policy_definition]
	p = sub, dom, obj, act
	
	[role_definition]
	g = _,_,_
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act`)
	if err != nil {
		log.Info().Msgf("fail to load model, error is: %v", err)
		return
	}
	enforcer, err := casbin.NewEnforcer(m, policy)
	if err != nil {
		log.Info().Msgf("init auth model error: %v", err)
		return
	}

	enforcer.EnableLog(true) // 开启权限认证日志
	// 加载数据库中的策略
	err = enforcer.LoadPolicy()
	if err != nil {
		log.Info().Msgf("loadPolicy error: %v", err)
		panic(err)
	}

	// 创建一个角色,并赋于权限
	// root 这个角色可以访问GET 方式访问 /api/v2/ping
	var ok bool
	ok, err = enforcer.AddPolicy("root", "xxx.com", "/api/v2/ping", "GET")
	if !ok {
		log.Info().Msg("policy is exist")
	}
	log.Info().Msgf("policy add success...")

	enforcer.AddRoleForUser("user", "root", "xxx.com")
	enforcer.AddRoleForUser("user1", "root1")
	// enforcer.DeleteUser("test")

	result, err := enforcer.Enforce("user", "xxx.com", "/api/v2/ping", "GET")
	if err != nil {
		log.Info().Msgf("fail to verify user, error is: %v", err)
		return
	}

	log.Info().Msgf("verify result is: %v", result)
}

// RBAC模式
func CasbinRbac() {

	engine, err := gorm.Open(sqlite.Open("debug/casbin/gorm_rbac.db"), &gorm.Config{})
	if err != nil {
		log.Info().Msgf("create db error: %v", err)
		return
	}

	// 字符串模式
	// policy := `
	// p, alice, data1, read
	// p, bob, data2, write
	// p, data2_admin, data2, read
	// p, data2_admin, data2, write
	// g, alice, data2_admin
	// `
	// policy := NewStringAdapter(policy)
	policy, err := gormadapter.NewAdapterByDB(engine)
	if err != nil {
		log.Info().Msgf("init auth error: %v", err)
		return
	}

	// enforcer, err := casbin.NewEnforcer("./auth_model.conf", policy)
	m, err := casbin_model.NewModelFromString(`[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[role_definition]
	g = _,_
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act`)
	if err != nil {
		log.Info().Msgf("fail to load model, error is: %v", err)
		return
	}
	enforcer, err := casbin.NewEnforcer(m, policy)
	if err != nil {
		log.Info().Msgf("init auth model error: %v", err)
		return
	}

	enforcer.EnableLog(true) // 开启权限认证日志
	// 加载数据库中的策略
	err = enforcer.LoadPolicy()
	if err != nil {
		log.Info().Msgf("loadPolicy error: %v", err)
		panic(err)
	}

	// 创建一个角色,并赋于权限
	// root 这个角色可以访问GET 方式访问 /api/v2/ping
	var ok bool
	ok, err = enforcer.AddPolicy("root", "/api/v2/ping", "GET")
	if !ok {
		log.Info().Msg("policy is exist")
	}
	log.Info().Msgf("policy add success...")

	enforcer.AddRoleForUser("user", "root")

	result, err := enforcer.Enforce("user", "/api/v2/ping", "GET")
	if err != nil {
		log.Info().Msgf("fail to verify user, error is: %v", err)
		return
	}

	log.Info().Msgf("verify result is: %v", result)
}

func CasbinAcl() {

	engine, err := gorm.Open(sqlite.Open("debug/casbin/gorm_acl.db"), &gorm.Config{})
	if err != nil {
		log.Info().Msgf("create db error: %v", err)
		return
	}

	policy, err := gormadapter.NewAdapterByDB(engine)
	if err != nil {
		log.Info().Msgf("init auth error: %v", err)
		return
	}

	enforcer, err := casbin.NewEnforcer("./debug/casbin/auth_model.conf", policy)
	if err != nil {
		log.Info().Msgf("init auth model error: %v", err)
		return
	}

	enforcer.EnableLog(true) // 开启权限认证日志
	// 加载数据库中的策略
	err = enforcer.LoadPolicy()
	if err != nil {
		log.Info().Msgf("loadPolicy error: %v", err)
		panic(err)
	}

	// 创建一个角色,并赋于权限
	// root 这个角色可以访问GET 方式访问 /api/v2/ping
	var ok bool
	ok, err = enforcer.AddPolicy("user", "/api/v2/ping", "GET")
	if !ok {
		log.Info().Msg("policy is exist")
	}
	log.Info().Msgf("policy add success...")

	result, err := enforcer.Enforce("user", "/api/v2/ping", "GET")
	if err != nil {
		log.Info().Msgf("fail to verify user, error is: %v", err)
		return
	}

	log.Info().Msgf("verify result is: %v", result)
}

func main() {
	// CasbinRbacWithDomain()
	// CasbinRbac()
	CasbinAcl()
}
