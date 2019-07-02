package remotehost

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
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

// NewClient creates a new SFTP client
func NewClient(source string) (*Client, error) {
	connection, err := CreateConnection(source)
	if err != nil {
		return nil, err
	}

	// create new SFTP client
	client, err := sftp.NewClient(connection)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create new SFTP client")
	}

	return &Client{client}, nil
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

// GetFile copies source file to destination file
func (c *Client) GetFile(srcFilePath string, dstFilePath string) (*int64, error) {
	srcFile, err := c.Open(srcFilePath)
	if err != nil {
		// int64(0) empty value to return in case of error
		return nil, err
	}

	defer srcFile.Close()
	os.MkdirAll(path.Dir(dstFilePath), 0755)

	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		return nil, err
	}

	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return nil, err
	}

	return &bytes, err
}

// Fetch2 retrieves a file
func Fetch2(hostname, src, dst string) error {
	client, err := NewClient(hostname)
	if err != nil {
		return err
	}

	defer client.Close()

	bytes, err := client.GetFile(src, dst)
	if err != nil {
		return errors.Wrap(err, "Cannot fetch file")
	}

	logrus.Printf("SFTP: %s:%s: %d bytes copied", hostname, src, bytes)
	return nil
}

// Fetch retrieves a file
func Fetch(hostname, src, dst string) error {
	client, err := NewClient(hostname)
	if err != nil {
		return err
	}

	defer client.Close()

	bytes, err := client.GetFile(src, dst)
	if err != nil {
		return errors.Wrap(err, "Cannot fetch file")
	}

	logrus.Printf("SFTP: %s:%s: %d bytes copied", hostname, src, bytes)
	return nil
}

// RunCMD executre cmd on remote host
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
// GetEnvVar get env var from remote host
func GetEnvVar(hostname, envVar string) (string, error) {
	session, err := NewSSHSession(hostname)
	if err != nil {
		return "", err
	}

	cmd := fmt.Sprintf("print $%s", envVar)
	output, err := session.Output(cmd)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
