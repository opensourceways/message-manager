/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package config define config interface.
package config

type configValidate interface {
	Validate() error
}

type configSetDefault interface {
	SetDefault()
}

type configItems interface {
	ConfigItems() []interface{}
}

// SetDefault sets default values for the provided configuration.
func SetDefault(cfg interface{}) {
	if f, ok := cfg.(configSetDefault); ok {
		f.SetDefault()
	}

	if f, ok := cfg.(configItems); ok {
		items := f.ConfigItems()

		for i := range items {
			SetDefault(items[i])
		}
	}
}

// Validate performs validation for the provided configuration.
func Validate(cfg interface{}) error {
	if f, ok := cfg.(configValidate); ok {
		if err := f.Validate(); err != nil {
			return err
		}
	}

	if f, ok := cfg.(configItems); ok {
		items := f.ConfigItems()

		for i := range items {
			if err := Validate(items[i]); err != nil {
				return err
			}
		}
	}

	return nil
}
