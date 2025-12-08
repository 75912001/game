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
	err = cfg.GCommon.Load()
	if err != nil {
		return errors.WithMessagef(err, "load common error %v", xruntime.Location())
	}
	err = cfg.GPetMgr.Load()
	if err != nil {
		return errors.WithMessagef(err, "load petMgr error %v", xruntime.Location())
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
	err = cfg.GCommon.Check()
	if err != nil {
		return errors.WithMessagef(err, "check common error %v", xruntime.Location())
	}
	err = cfg.GPetMgr.Check()
	if err != nil {
		return errors.WithMessagef(err, "check petMgr error %v", xruntime.Location())
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
	err = cfg.GCommon.Assemble()
	if err != nil {
		return errors.WithMessagef(err, "Assemble common error %v", xruntime.Location())
	}
	err = cfg.GPetMgr.Assemble()
	if err != nil {
		return errors.WithMessagef(err, "Assemble petMgr error %v", xruntime.Location())
	}
	return nil
}
