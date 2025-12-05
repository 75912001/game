package main

import "saClient/src/cfg"

// LoadCfg 加载配置
func LoadCfg() error {
	err := cfg.GRoleMgr.Load()
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
	return nil
}
