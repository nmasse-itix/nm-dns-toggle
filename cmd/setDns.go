/*
Copyright © 2022 Nicolas MASSE

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/spf13/cobra"
)

var servers []string

// setDnsCmd represents the setDns command
var setDnsCmd = &cobra.Command{
	Use:   "set-dns",
	Short: "Sets upstream DNS of dnsmasq",
	Long: `While initializing a Wireguard connection, temporarily sets the upstream
servers used by dnsmasq.`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := dbus.ConnectSystemBus()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to connect to system bus:", err)
			os.Exit(1)
		}
		defer conn.Close()

		fmt.Printf("Setting the dnsmasq upstream servers to: %s...\n", strings.Join(servers, ", "))

		obj := conn.Object("org.freedesktop.NetworkManager.dnsmasq", "/uk/org/thekelleys/dnsmasq")
		err = obj.Call("uk.org.thekelleys.dnsmasq.SetDomainServers", 0, servers).Store()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to call uk.org.thekelleys.dnsmasq.SetDomainServers on org.freedesktop.NetworkManager.dnsmasq", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(setDnsCmd)
	setDnsCmd.Flags().StringSliceVarP(&servers, "server", "s", []string{}, "Upstream servers to set")
}
