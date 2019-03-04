package cmd

// All commands should implement this interfaces
type Commander interface {
	Execute(args []string) error
}
