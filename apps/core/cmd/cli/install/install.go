package install

import (
	"github.com/spf13/cobra"
)

const (
	DefaultUserName = "web master"
)

var (
	cfgFile     string
	userLicense string
	username    string
	email       string
	password    string

	InstallCmd = &cobra.Command{
		Use:   "install",
		Short: "Install kova on this machine or a remote machine",
		Long:  `Install kova on this machine or a remote machine, use's SSH for remote machine `,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func Execute() error {
	return InstallCmd.Execute()
}

func init() {
	InstallCmd.Flags().StringVarP(&username, "name", "n", DefaultUserName, "Username of the admin user, defaults to web master")
	InstallCmd.Flags().StringVarP(&email, "email", "e", "", "Email of the admin user (required)")
	InstallCmd.Flags().StringVarP(&password, "password", "p", "", "Password of the admin user (required)")
	InstallCmd.MarkFlagRequired("password")
	InstallCmd.MarkFlagRequired("email")
}
