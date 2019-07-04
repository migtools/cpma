package remotehost

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/fusor/cpma/pkg/env"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

var instances struct {
	client   *ssh.Client
	sshError error
}

var once struct {
	client sync.Once
}

// Client Wrapper around sftp.Client
type Client struct {
	*sftp.Client
}

// CreateConnection create ssh connection
func CreateConnection(source string) (*ssh.Client, error) {
	key, err := ioutil.ReadFile(env.Config().GetString("SSHPrivateKey"))
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to read private key: %s\n", key)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse private key")
	}

	knownHostsFile := filepath.Join(env.Config().GetString("home"), ".ssh", "known_hosts")

	var hostKeyCallback ssh.HostKeyCallback
	if env.Config().GetBool("InsecureHostKey") {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	} else {
		hostKeyCallback, err = kh.New(knownHostsFile)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to get hostkey in %s\n", knownHostsFile)
		}
	}

	sshConfig := &ssh.ClientConfig{
		User: env.Config().GetString("SSHLogin"),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,

		Timeout: 10 * time.Second,
	}

	port := 22
	if p := env.Config().GetString("SSHPort"); p != "" {
		port, err = strconv.Atoi(p)
		if err != nil || port < 1 || port > 65535 {
			return nil, errors.Wrapf(err, "Port number %s is wrong\n", p)
		}
	}

	addr := fmt.Sprintf("%s:%d", source, port)
	logrus.Debugf("Connecting to %s", addr)

	connection, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot connect to %s\n", addr)
	}

	return connection, nil
}

// NewSSHSession Start new ssh session
func NewSSHSession(source string) (*ssh.Session, error) {
	// Init client only on first call
	once.client.Do(func() {
		connection, err := CreateConnection(source)
		instances.client = connection
		instances.sshError = err
	})

	if instances.sshError != nil {
		return nil, instances.sshError
	}

	session, err := instances.client.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "Cannot start session")
	}

	return session, nil
}

// RunCMD execute cmd on remote host
func RunCMD(hostname, cmd string) (string, error) {
	session, err := NewSSHSession(hostname)
	if err != nil {
		return "", err
	}

	output, err := session.Output(cmd)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
