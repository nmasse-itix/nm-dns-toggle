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

// revertCmd represents the revert command
var revertCmd = &cobra.Command{
	Use:   "revert",
	Short: "Reverts upstream DNS of dnsmasq",
	Long: `Reverts upstream DNS servers to their original values when the wireguard
connection is closed.`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := dbus.ConnectSystemBus()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to connect to system bus:", err)
			os.Exit(1)
		}
		defer conn.Close()

		var entries []map[string]dbus.Variant
		obj := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager/DnsManager")
		err = obj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.NetworkManager.DnsManager", "Configuration").Store(&entries)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to fetch org.freedesktop.NetworkManager.DnsManager.Configuration property on org.freedesktop.NetworkManager:", err)
			os.Exit(1)
		}

		for _, entry := range entries {
			for k, v := range entry {
				if k == "nameservers" {
					if nameservers, ok := v.Value().([]string); ok {
						servers = append(servers, nameservers...)
					}
				}
			}
		}
		fmt.Printf("Setting the dnsmasq upstream servers to: %s...\n", strings.Join(servers, ", "))
		obj = conn.Object("org.freedesktop.NetworkManager.dnsmasq", "/uk/org/thekelleys/dnsmasq")
		err = obj.Call("uk.org.thekelleys.dnsmasq.SetDomainServers", 0, servers).Store()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to call uk.org.thekelleys.dnsmasq.SetDomainServers on org.freedesktop.NetworkManager.dnsmasq", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(revertCmd)
}
