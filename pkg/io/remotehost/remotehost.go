package remotehost

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fusor/cpma/pkg/env"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

// Client Wrapper around sftp.Client
type Client struct {
	*sftp.Client
}

// CreateConnection create ssh connection
func CreateConnection(source string) (*ssh.Client, error) {
	sshCreds := env.Config().GetStringMapString("SSHCreds")

	key, err := ioutil.ReadFile(sshCreds["privatekey"])
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to read private key: %s\n", sshCreds["privatekey"])
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
		User: sshCreds["login"],
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,

		Timeout: 10 * time.Second,
	}

	port := 22
	if p := sshCreds["port"]; p != "" {
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
	connection, err := CreateConnection(source)
	if err != nil {
		return nil, err
	}

	session, err := connection.NewSession()
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
