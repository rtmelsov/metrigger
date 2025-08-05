// Package helpers
package helpers

import (
	"github.com/rtmelsov/metrigger/internal/models"
	"strconv"
)

func ClientFileConfigParser(flags *models.AgentFlags, confs *models.ClientFileConfig) error {
	if flags.Addr == "" {
		flags.Addr = confs.Address
	}
	if flags.ReportInterval == 0 {
		n, err := strconv.Atoi(string(confs.ReportInterval[0]))
		if err != nil {
			return err
		}
		flags.ReportInterval = n
	}
	if flags.PollInterval == 0 {
		n, err := strconv.Atoi(string(confs.PollInterval[0]))
		if err != nil {
			return err
		}
		flags.PollInterval = n
	}
	if flags.CryptoRate == "" {
		flags.CryptoRate = confs.CryptoKey
	}
	if flags.TrustedSubnet == "" {
		flags.TrustedSubnet = confs.TrustedSubnet
	}
	return nil
}

func ServerFileConfigParser(flags *models.ServerFlagsType, confs *models.ServerFileConfig) error {
	if !flags.Restore {
		flags.Restore = confs.Restore
	}

	if flags.Addr == "" {
		flags.Addr = confs.Address
	}
	if flags.CryptoRate == "" {
		flags.CryptoRate = confs.CryptoKey
	}
	if flags.DataBaseDsn == "" {
		flags.DataBaseDsn = confs.DataBaseDsn
	}
	if flags.FileStoragePath == "" {
		flags.FileStoragePath = confs.StoreFile
	}
	if flags.TrustedSubnet == "" {
		flags.TrustedSubnet = confs.TrustedSubnet
	}

	if flags.StoreInterval == 0 {
		n, err := strconv.Atoi(string(confs.StoreInterval[0]))
		if err != nil {
			return err
		}
		flags.StoreInterval = n
	}
	return nil
}
