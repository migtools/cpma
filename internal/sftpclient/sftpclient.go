package sftpclient

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/fusor/cpma/internal/config"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

type Info struct {
	Login      string `mapstructure:"Login"`
	PrivateKey string `mapstructure:"PrivateKey"`
}

type Client struct {
	*sftp.Client
}

// NewClient creates a new SFTP client
func (c *Info) NewClient(hostname string) (*Client, error) {
	key, err := ioutil.ReadFile(c.PrivateKey)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	hostKeyCallback, err := kh.New(filepath.Join(config.Config().GetString("home"), ".ssh", "known_hosts"))
	if err != nil {
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User: c.Login,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,

		Timeout: 10 * time.Second,
	}

	// TODO: accept custom port
	addr := fmt.Sprintf("%s:22", hostname)
	connection, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, err
	}

	// create new SFTP client
	client, err := sftp.NewClient(connection)
	if err != nil {
		return nil, err
	}

	return &Client{client}, err
}

// GetFile copies source file to destination file
func (c *Client) GetFile(srcFilePath string, dstFilePath string) {
	srcFile, err := c.Open(srcFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	os.MkdirAll(path.Dir(dstFilePath), 0755)
	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("File %s: %d bytes copied\n", srcFilePath, bytes)
}
