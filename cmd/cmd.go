package cmd

import "shopaholic/store/service"

// All commands should implement this interfaces
type Commander interface {
	Execute(args []string) error
	SetCommon(commonOpts CommonOpts)
}

type CommonOpts struct {
	Currency string
	Store    service.DataStore
}

func (c *CommonOpts) SetCommon(commonOpts CommonOpts) {
	c.Currency = commonOpts.Currency
	c.Store = commonOpts.Store
}
