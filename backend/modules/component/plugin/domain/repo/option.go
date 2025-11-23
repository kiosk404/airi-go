package repo

import (
	"github.com/kiosk404/airi-go/backend/modules/component/plugin/infra/repo"
)

type PluginSelectedOptions func(*repo.PluginSelectedOption)

func WithPluginID() PluginSelectedOptions {
	return func(opts *repo.PluginSelectedOption) {
		opts.PluginID = true
	}
}

func WithPluginOpenapiDoc() PluginSelectedOptions {
	return func(opts *repo.PluginSelectedOption) {
		opts.OpenapiDoc = true
	}
}

func WithPluginManifest() PluginSelectedOptions {
	return func(opts *repo.PluginSelectedOption) {
		opts.Manifest = true
	}
}

func WithPluginIconURI() PluginSelectedOptions {
	return func(opts *repo.PluginSelectedOption) {
		opts.IconURI = true
	}
}

func WithPluginVersion() PluginSelectedOptions {
	return func(opts *repo.PluginSelectedOption) {
		opts.Version = true
	}
}

type ToolSelectedOptions func(option *repo.ToolSelectedOption)

func WithToolID() ToolSelectedOptions {
	return func(opts *repo.ToolSelectedOption) {
		opts.ToolID = true
	}
}

func WithToolMethod() ToolSelectedOptions {
	return func(opts *repo.ToolSelectedOption) {
		opts.ToolMethod = true
	}
}

func WithToolSubURL() ToolSelectedOptions {
	return func(opts *repo.ToolSelectedOption) {
		opts.ToolSubURL = true
	}
}

func WithToolActivatedStatus() ToolSelectedOptions {
	return func(opts *repo.ToolSelectedOption) {
		opts.ActivatedStatus = true
	}
}

func WithToolDebugStatus() ToolSelectedOptions {
	return func(opts *repo.ToolSelectedOption) {
		opts.DebugStatus = true
	}
}
