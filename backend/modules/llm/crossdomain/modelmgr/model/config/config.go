package config

import (
	"fmt"
)

type BasicConfiguration struct {
	CodeRunnerType CodeRunnerType `json:"code_runner_type"`
	SandboxConfig  *SandboxConfig `json:"sandbox_config,omitempty"`
}

type CodeRunnerType int64

const (
	CoderunnerTypeLocal   CodeRunnerType = 0
	CoderunnerTypeSandbox CodeRunnerType = 1
)

func (p CodeRunnerType) String() string {
	switch p {
	case CoderunnerTypeLocal:
		return "Local"
	case CoderunnerTypeSandbox:
		return "Sandbox"
	}
	return "<UNSET>"
}

func CodeRunnerTypeFromString(s string) (CodeRunnerType, error) {
	switch s {
	case "Local":
		return CoderunnerTypeLocal, nil
	case "Sandbox":
		return CoderunnerTypeSandbox, nil
	}
	return CodeRunnerType(0), fmt.Errorf("not a valid CodeRunnerType string")
}

type SandboxConfig struct {
	AllowEnv       string  `json:"allow_env"`
	AllowRead      string  `json:"allow_read"`
	AllowWrite     string  `json:"allow_write"`
	AllowRun       string  `json:"allow_run"`
	AllowNet       string  `json:"allow_net"`
	AllowFfi       string  `json:"allow_ffi"`
	NodeModulesDir string  `json:"node_modules_dir"`
	TimeoutSeconds float64 `json:"timeout_seconds"`
	MemoryLimitMb  int64   `json:"memory_limit_mb"`
}

func NewSandboxConfig() *SandboxConfig {
	return &SandboxConfig{}
}

func (p *SandboxConfig) InitDefault() {
}

func (p *SandboxConfig) GetAllowEnv() (v string) {
	return p.AllowEnv
}

func (p *SandboxConfig) GetAllowRead() (v string) {
	return p.AllowRead
}

func (p *SandboxConfig) GetAllowWrite() (v string) {
	return p.AllowWrite
}

func (p *SandboxConfig) GetAllowRun() (v string) {
	return p.AllowRun
}

func (p *SandboxConfig) GetAllowNet() (v string) {
	return p.AllowNet
}

func (p *SandboxConfig) GetAllowFfi() (v string) {
	return p.AllowFfi
}

func (p *SandboxConfig) GetNodeModulesDir() (v string) {
	return p.NodeModulesDir
}

func (p *SandboxConfig) GetTimeoutSeconds() (v float64) {
	return p.TimeoutSeconds
}

func (p *SandboxConfig) GetMemoryLimitMb() (v int64) {
	return p.MemoryLimitMb
}
func (p *SandboxConfig) SetAllowEnv(val string) {
	p.AllowEnv = val
}
func (p *SandboxConfig) SetAllowRead(val string) {
	p.AllowRead = val
}
func (p *SandboxConfig) SetAllowWrite(val string) {
	p.AllowWrite = val
}
func (p *SandboxConfig) SetAllowRun(val string) {
	p.AllowRun = val
}
func (p *SandboxConfig) SetAllowNet(val string) {
	p.AllowNet = val
}
func (p *SandboxConfig) SetAllowFfi(val string) {
	p.AllowFfi = val
}
func (p *SandboxConfig) SetNodeModulesDir(val string) {
	p.NodeModulesDir = val
}
func (p *SandboxConfig) SetTimeoutSeconds(val float64) {
	p.TimeoutSeconds = val
}
func (p *SandboxConfig) SetMemoryLimitMb(val int64) {
	p.MemoryLimitMb = val
}

func (p *SandboxConfig) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SandboxConfig(%+v)", *p)
}
