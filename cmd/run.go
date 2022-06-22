/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yndd/ztp-dhcp/pkg/backend/k8s"
	_ "github.com/yndd/ztp-dhcp/pkg/devices/all"
	"github.com/yndd/ztp-dhcp/pkg/dhcp"
)

var (
	dhcpv4_port int
	//dhcpv6_port int
	ifname    string // specific interface to start dhcp server on
	leaseTime uint32
	// ip address the dhcp server is bound to
	listenip string
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the ZTP DHCP-Server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// init the k8s backend
		//backend := k8s.NewZtpK8sBackend()

		// prepare the static Backend
		backend := k8s.NewZtpK8sBackend("")

		// setup the server
		ztpserver := dhcp.NewZtpServer(backend, &dhcp.ZtpSettings{LeaseTime: leaseTime})
		// execute the server
		ztpserver.StartDHCPServer(dhcpv4_port, ifname, listenip)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// TODO: allow for viper environment variables here
	// viper.SetEnvPrefix("NDD")
	// viper.BindEnv("dhcpv4_port")

	// viper.BindPFlag("dhcpv4_port", runCmd.Flags().Lookup("dhcpv4-port"))

	runCmd.Flags().IntVar(&dhcpv4_port, "dhcpv4-port", 67, "The port to bind the dhcpv4 server to.")
	//runCmd.Flags().IntVar(&dhcpv6_port, "dhcpv6-port", 567, "The port to bind the dhcpv6 server to.")
	runCmd.Flags().StringVar(&ifname, "interface", "", "Define the interface to bind the DHCP server to. If left empty [default] the server is not bound to a specific interface.")
	runCmd.Flags().StringVar(&listenip, "listen-ip", "0.0.0.0", "IP-Address the DHCP server is bound to")
	runCmd.Flags().Uint32Var(&leaseTime, "lease-time", 3600, "The lease time in seconds.")
}
