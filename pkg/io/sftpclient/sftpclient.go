package sftpclient

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fusor/cpma/pkg/env"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

// Client Wrapper around sftp.Client
type Client struct {
	*sftp.Client
}

// NewClient creates a new SFTP client
func NewClient(source string) (Client, error) {
	sshCreds := env.Config().GetStringMapString("SSHCreds")

	key, err := ioutil.ReadFile(sshCreds["privatekey"])
	if err != nil {
		logrus.Errorf("Unable to read private key: %s", sshCreds["privatekey"])
		return Client{}, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		logrus.Error("Unable to parse private key")
		return Client{}, err
	}

	knownHostsFile := filepath.Join(env.Config().GetString("home"), ".ssh", "known_hosts")

	hostKeyCallback, err := kh.New(knownHostsFile)
	if err != nil {
		logrus.Errorf("Unable to get hostkey in %s", knownHostsFile)
		return Client{}, err
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
			return Client{}, errors.New("Port number " + p + " is wrong.")
		}
	}

	addr := fmt.Sprintf("%s:%d", source, port)
	logrus.Debug("Connecting to", addr)

	connection, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		logrus.Errorf("Cannot connect to %s", addr)
		return Client{}, err
	}

	// create new SFTP client
	client, err := sftp.NewClient(connection)
	if err != nil {
		logrus.Error("Unable to create new SFTP client")
		return Client{}, err
	}

	return Client{client}, nil
}

// GetFile copies source file to destination file
func (c *Client) GetFile(srcFilePath string, dstFilePath string) (int64, error) {
	srcFile, err := c.Open(srcFilePath)
	if err != nil {
		// int64(0) empty value to return in case of error
		return int64(0), err
	}

	defer srcFile.Close()
	os.MkdirAll(path.Dir(dstFilePath), 0755)

	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		return int64(0), err
	}

	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return int64(0), err
	}

	return bytes, err
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
		return err
	}

	logrus.Printf("SFTP: %s:%s: %d bytes copied", hostname, src, bytes)
	return nil
}
