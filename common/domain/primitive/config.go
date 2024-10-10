/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package primitive define the domain of object
package primitive

import (
	"regexp"

	"k8s.io/apimachinery/pkg/util/sets"
)

type Config struct {
	MSD            MSDConfig     `json:"msd"`
	File           FileConfig    `json:"file"`
	Account        AccountConfig `json:"account"`
	RandomIdLength int           `json:"random_id_length"`
}

// SetDefault sets default values for Config if they are not provided.
func (cfg *Config) SetDefault() {
	if cfg.RandomIdLength <= 0 {
		cfg.RandomIdLength = 24
	}
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.MSD,
		&cfg.File,
		&cfg.Account,
	}
}

// MSDConfig represents the configuration for MSD.
type MSDConfig struct {
	NameRegexp    string `json:"msd_name_regexp"          required:"true"`
	MinNameLength int    `json:"msd_name_min_length"      required:"true"`
	MaxNameLength int    `json:"msd_name_max_length"      required:"true"`

	nameRegexp *regexp.Regexp
}

// Validate validates values for MSDConfig whether they are valid.
func (cfg *MSDConfig) Validate() (err error) {
	cfg.nameRegexp, err = regexp.Compile(cfg.NameRegexp)

	return
}

// FileConfig represents the configuration for file.
type FileConfig struct {
	FileRefRegexp     string `json:"file_ref_regexp"          required:"true"`
	FileNameRegexp    string `json:"file_name_regexp"         required:"true"`
	FileRefMinLength  int    `json:"file_ref_min_length"      required:"true"`
	FileRefMaxLength  int    `json:"file_ref_max_length"      required:"true"`
	FilePathMaxLength int    `json:"file_path_max_length"     required:"true"`

	fileRefRegexp  *regexp.Regexp
	fileNameRegexp *regexp.Regexp
}

// Validate validates values for fileConfig whether they are valid.
func (cfg *FileConfig) Validate() (err error) {
	if cfg.fileRefRegexp, err = regexp.Compile(cfg.FileRefRegexp); err != nil {
		return err
	}

	if cfg.fileNameRegexp, err = regexp.Compile(cfg.FileNameRegexp); err != nil {
		return err
	}

	return nil
}

// AccountConfig represents the configuration for Account.
type AccountConfig struct {
	NameRegexp       string   `json:"account_name_regexp"         required:"true"`
	MinAccountLength int      `json:"account_name_min_length"     required:"true"`
	MaxAccountLength int      `json:"account_name_max_length"     required:"true"`
	ReservedAccounts []string `json:"reserved_accounts"           required:"true"`

	nameRegexp       *regexp.Regexp
	reservedAccounts sets.Set[string]
}

// Validate validates values for AccountConfig whether they are valid.
func (cfg *AccountConfig) Validate() (err error) {
	if cfg.nameRegexp, err = regexp.Compile(cfg.NameRegexp); err != nil {
		return err
	}

	if len(cfg.ReservedAccounts) > 0 {
		cfg.reservedAccounts = sets.New[string]()
		cfg.reservedAccounts.Insert(cfg.ReservedAccounts...)
	}

	return nil
}
