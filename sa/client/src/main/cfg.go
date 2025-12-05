package main

import "saClient/src/cfg"

// LoadCfg 加载配置
func LoadCfg() error {
	err := cfg.GRoleMgr.Load()
	if err != nil {
		return err
	}
	err = cfg.GRoleBaseMgr.Load()
	if err != nil {
		return err
	}
	return nil
}

// CheckCfg 检查配置
func CheckCfg() error {
	err := cfg.GRoleMgr.Check()
	if err != nil {
		return err
	}
	err = cfg.GRoleBaseMgr.Check()
	if err != nil {
		return err
	}
	return nil
}

// AssembleCfg 装配配置
func AssembleCfg() error {
	err := cfg.GRoleMgr.Assemble()
	if err != nil {
		return err
	}
	err = cfg.GRoleBaseMgr.Assemble()
	if err != nil {
		return err
	}
	return nil
}
