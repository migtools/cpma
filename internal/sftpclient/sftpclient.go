package sftpclient

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fusor/cpma/env"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

// Client Wrapper around sftp.Client
type Client struct {
	*sftp.Client
}

func NewClient(source string) Client {
	sshCreds := env.Config().GetStringMapString("SSHCreds")

	key, err := ioutil.ReadFile(sshCreds["privatekey"])
	if err != nil {
		logrus.WithError(err).Fatalf("Unable to read private key: %s", sshCreds["privatekey"])
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		logrus.Fatalf("Unable to parse private key: %v", err)
	}

	knownHostsFile := filepath.Join(env.Config().GetString("home"), ".ssh", "known_hosts")
	hostKeyCallback, err := kh.New(knownHostsFile)
	if err != nil {
		logrus.WithError(err).Fatalf("Unable to get hostkey in %s", knownHostsFile)
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
			logrus.Fatalf("Port number (%s) is wrong.", p)
		}
	}

	addr := fmt.Sprintf("%s:%d", source, port)
	logrus.Debug("Connecting to", addr)

	connection, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		logrus.WithError(err).Fatalf("Cannot connect to %s", addr)
	}

	// create new SFTP client
	client, err := sftp.NewClient(connection)
	if err != nil {
		logrus.Fatal(err)
	}

	return Client{client}
}

// GetFile copies source file to destination file
func (c *Client) GetFile(srcFilePath string, dstFilePath string) (int64, error) {
	srcFile, err := c.Open(srcFilePath)
	if err != nil {
		logrus.Fatal(err)
	}
	defer srcFile.Close()

	os.MkdirAll(path.Dir(dstFilePath), 0755)
	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		logrus.Fatal(err)
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		logrus.Fatal(err)
	}
	return bytes, err
}

// Fetch retrieves a file
func Fetch(hostname, src, dst string) {
	client := NewClient(hostname)
	defer client.Close()

	bytes, err := client.GetFile(src, dst)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Printf("File %s:%s: %d bytes copied", hostname, src, bytes)
}
