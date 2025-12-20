package main

import (
	"saClient/src/cfg"
	"saClient/src/res"

	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// LoadCfg 加载配置
func LoadCfg() error {
	err := cfg.GRoleMgr.Load()
	if err != nil {
		return errors.WithMessagef(err, "load role error %v", xruntime.Location())
	}
	err = res.GRoleMgr.Load()
	if err != nil {
		return errors.WithMessagef(err, "load res role error %v", xruntime.Location())
	}
	err = res.GPetMgr.Load()
	if err != nil {
		return errors.WithMessagef(err, "load res pet error %v", xruntime.Location())
	}
	err = cfg.GCommon.Load()
	if err != nil {
		return errors.WithMessagef(err, "load common error %v", xruntime.Location())
	}
	err = cfg.GPetMgr.Load()
	if err != nil {
		return errors.WithMessagef(err, "load petMgr error %v", xruntime.Location())
	}
	err = cfg.GSceneTransitionMgr.Load()
	if err != nil {
		return errors.WithMessagef(err, "load transitionMgr error %v", xruntime.Location())
	}
	err = cfg.GMapMgr.Load()
	if err != nil {
		return errors.WithMessagef(err, "load mapMgr error %v", xruntime.Location())
	}
	err = cfg.GPortalMgr.Load()
	if err != nil {
		return errors.WithMessagef(err, "load teleportMgr error %v", xruntime.Location())
	}
	// 加载 Tiled 地图资源
	err = cfg.GTiledMapMgr.Load()
	if err != nil {
		return errors.WithMessagef(err, "load res tiledMapMgr error %v", xruntime.Location())
	}
	// 加载 enemy.group 配置
	err = cfg.GEnemyGroupMgr.Load()
	if err != nil {
		return errors.WithMessagef(err, "load enemy group error %v", xruntime.Location())
	}
	return nil
}

// CheckCfg 检查配置
func CheckCfg() error {
	err := cfg.GRoleMgr.Check()
	if err != nil {
		return errors.WithMessagef(err, "check role error %v", xruntime.Location())
	}
	err = res.GRoleMgr.Check()
	if err != nil {
		return errors.WithMessagef(err, "check res role error %v", xruntime.Location())
	}
	err = res.GPetMgr.Check()
	if err != nil {
		return errors.WithMessagef(err, "check res pet error %v", xruntime.Location())
	}
	err = cfg.GCommon.Check()
	if err != nil {
		return errors.WithMessagef(err, "check common error %v", xruntime.Location())
	}
	err = cfg.GPetMgr.Check()
	if err != nil {
		return errors.WithMessagef(err, "check petMgr error %v", xruntime.Location())
	}
	err = cfg.GSceneTransitionMgr.Check()
	if err != nil {
		return errors.WithMessagef(err, "check sceneTransitionMgr error %v", xruntime.Location())
	}
	err = cfg.GMapMgr.Check()
	if err != nil {
		return errors.WithMessagef(err, "check mapMgr error %v", xruntime.Location())
	}
	err = cfg.GPortalMgr.Check()
	if err != nil {
		return errors.WithMessagef(err, "check teleportMgr error %v", xruntime.Location())
	}
	// 检查 Tiled 地图资源
	err = cfg.GTiledMapMgr.Check()
	if err != nil {
		return errors.WithMessagef(err, "check res tiledMapMgr error %v", xruntime.Location())
	}
	// 检查 enemy.group 配置
	err = cfg.GEnemyGroupMgr.Check()
	if err != nil {
		return errors.WithMessagef(err, "check enemy group error %v", xruntime.Location())
	}
	return nil
}

// AssembleCfg 装配配置
func AssembleCfg() error {
	err := cfg.GRoleMgr.Assemble()
	if err != nil {
		return errors.WithMessagef(err, "Assemble role error %v", xruntime.Location())
	}
	err = res.GRoleMgr.Assemble()
	if err != nil {
		return errors.WithMessagef(err, "Assemble res role error %v", xruntime.Location())
	}
	err = res.GPetMgr.Assemble()
	if err != nil {
		return errors.WithMessagef(err, "Assemble res pet error %v", xruntime.Location())
	}
	err = cfg.GCommon.Assemble()
	if err != nil {
		return errors.WithMessagef(err, "Assemble common error %v", xruntime.Location())
	}
	err = cfg.GPetMgr.Assemble()
	if err != nil {
		return errors.WithMessagef(err, "Assemble petMgr error %v", xruntime.Location())
	}
	err = cfg.GSceneTransitionMgr.Assemble()
	if err != nil {
		return errors.WithMessagef(err, "Assemble sceneTransitionMgr error %v", xruntime.Location())
	}
	err = cfg.GMapMgr.Assemble()
	if err != nil {
		return errors.WithMessagef(err, "Assemble mapMgr error %v", xruntime.Location())
	}
	err = cfg.GPortalMgr.Assemble()
	if err != nil {
		return errors.WithMessagef(err, "Assemble teleportMgr error %v", xruntime.Location())
	}
	// 装配 Tiled 地图资源
	err = cfg.GTiledMapMgr.Assemble()
	if err != nil {
		return errors.WithMessagef(err, "Assemble res tiledMapMgr error %v", xruntime.Location())
	}
	// 装配 enemy.group 配置
	err = cfg.GEnemyGroupMgr.Assemble()
	if err != nil {
		return errors.WithMessagef(err, "Assemble enemy group error %v", xruntime.Location())
	}
	return nil
}
