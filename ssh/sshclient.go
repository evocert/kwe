package ssh

import (
	"fmt"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

/*type SSHConnConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}*/

type SSHConnection struct {
	host     string
	user     string
	password string
	sshclnts map[*sshclient]*sshclient
}

type sshclient struct {
	sshcn   *SSHConnection
	sshclnt *ssh.Client
}

func NewSSHConnection(host string, username string, password string) (sshcn *SSHConnection) {
	sshcn = &SSHConnection{host: host, sshclnts: map[*sshclient]*sshclient{}}
	return
}

func (sshcn *SSHConnection) Connect() (clnt *sshclient, err error) {
	if sshclnt, err := newSSHClient(sshcn); err == nil && sshclnt != nil {
		clnt = &sshclient{sshcn: sshcn, sshclnt: sshclnt}
	}
	return
}

func newSSHClient(sshcn *SSHConnection) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: sshcn.user,
		Auth: []ssh.AuthMethod{SSHAgent()},
	}

	if sshcn.password != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.Password(sshcn.password))
	}

	port := "22"
	host := sshcn.host
	if strings.HasPrefix(host, "[") && strings.Contains(host, "]:") {
		if !strings.HasSuffix(host, "]:") {
			port = host[strings.LastIndex(host, "]:")+1:]
		}
		host = host[len("["):strings.LastIndex(host, "]:")]
	} else if strings.Contains(host, ":") {
		if !strings.HasSuffix(host, ":") {
			port = host[strings.LastIndex(host, ":")+1:]
		}
		host = host[:strings.LastIndex(host, ":")]
	}

	if homeDir, err := os.UserHomeDir(); err == nil {
		if hostKeyCallback, err := knownhosts.New(fmt.Sprintf("%s/.ssh/known_hosts", homeDir)); err == nil {
			sshConfig.HostKeyCallback = hostKeyCallback
		}
	}
	return ssh.Dial("tcp", net.JoinHostPort(host, port), sshConfig)
}

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}
