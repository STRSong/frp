package basic

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatedier/frp/test/e2e/framework"

	. "github.com/onsi/ginkgo"
)

const (
	ConfigValidStr = "syntax is ok"
)

var _ = Describe("[Feature: Cmd]", func() {
	f := framework.NewDefaultFramework()

	Describe("Verify", func() {
		It("frps valid", func() {
			path := f.GenerateConfigFile(`
			[common]
			bind_addr = 0.0.0.0
			bind_port = 7000
			`)
			_, output, err := f.RunFrps([]string{"verify", "-c", path})
			framework.ExpectNoError(err)
			framework.ExpectTrue(strings.Contains(output, ConfigValidStr), "output: %s", output)
		})
		It("frps invalid", func() {
			path := f.GenerateConfigFile(`
			[common]
			bind_addr = 0.0.0.0
			bind_port = 70000
			`)
			_, output, err := f.RunFrps([]string{"verify", "-c", path})
			framework.ExpectNoError(err)
			framework.ExpectTrue(!strings.Contains(output, ConfigValidStr), "output: %s", output)
		})
		It("frpc valid", func() {
			path := f.GenerateConfigFile(`
			[common]
			server_addr = 0.0.0.0
			server_port = 7000
			`)
			_, output, err := f.RunFrpc([]string{"verify", "-c", path})
			framework.ExpectNoError(err)
			framework.ExpectTrue(strings.Contains(output, ConfigValidStr), "output: %s", output)
		})
		It("frpc invalid", func() {
			path := f.GenerateConfigFile(`
			[common]
			server_addr = 0.0.0.0
			server_port = 7000
			protocol = invalid
			`)
			_, output, err := f.RunFrpc([]string{"verify", "-c", path})
			framework.ExpectNoError(err)
			framework.ExpectTrue(!strings.Contains(output, ConfigValidStr), "output: %s", output)
		})
	})

	Describe("Single proxy", func() {
		It("TCP", func() {
			serverPort := f.AllocPort()
			_, _, err := f.RunFrps([]string{"-t", "123", "-p", strconv.Itoa(serverPort)})
			framework.ExpectNoError(err)

			localPort := f.PortByName(framework.TCPEchoServerPort)
			remotePort := f.AllocPort()
			f.RunFrpc([]string{"tcp", "-s", fmt.Sprintf("127.0.0.1:%d", serverPort), "-t", "123", "-u", "test",
				"-l", strconv.Itoa(localPort), "-r", strconv.Itoa(remotePort), "-n", "tcp_test"})

			framework.NewRequestExpect(f).Port(remotePort).Ensure()
		})

		It("UDP", func() {
			serverPort := f.AllocPort()
			_, _, err := f.RunFrps([]string{"-t", "123", "-p", strconv.Itoa(serverPort)})
			framework.ExpectNoError(err)

			localPort := f.PortByName(framework.UDPEchoServerPort)
			remotePort := f.AllocPort()
			f.RunFrpc([]string{"udp", "-s", fmt.Sprintf("127.0.0.1:%d", serverPort), "-t", "123", "-u", "test",
				"-l", strconv.Itoa(localPort), "-r", strconv.Itoa(remotePort), "-n", "udp_test"})

			framework.NewRequestExpect(f).RequestModify(framework.SetRequestProtocol("udp")).
				Port(remotePort).Ensure()
		})
	})
})