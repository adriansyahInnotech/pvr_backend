package config

import (
	"fmt"
	"os"
	"path"
	"sync"

	"pvr_backend/db"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var (
	enforcer *casbin.Enforcer
	once     sync.Once
)

func InitCasbin() {
	once.Do(func() {
		adapter, err := gormadapter.NewAdapterByDB(db.GetDB())
		if err != nil {
			panic(fmt.Errorf("create casbin adapter: %w", err))
		}

		cwd, _ := os.Getwd()
		modelPath := path.Join(cwd, "config", "rbac_model.conf")

		enf, err := casbin.NewEnforcer(modelPath, adapter)
		if err != nil {
			panic(fmt.Errorf("create casbin enforcer: %w", err))
		}

		if err := enf.LoadPolicy(); err != nil {
			panic(fmt.Errorf("load policy: %w", err))
		}

		fmt.Println("✅ Casbin policy loaded")
		enforcer = enf
	})
}

func GetEnforcer() *casbin.Enforcer {
	if enforcer == nil {
		panic("Casbin not initialized. Call InitCasbin() first")
	}
	return enforcer
}
