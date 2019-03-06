package sftpclient

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/pkg/sftp"
)

func GetFile(host string, user string, keyfile string, srcFilePath string, dstFilePath string) {
	// get host public key
	hostKey := getHostKey(host)

	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},

		// verify host public key
		HostKeyCallback: ssh.FixedHostKey(hostKey),

		Timeout: 10 * time.Second,
	}

	connection, err := ssh.Dial("tcp", host+":"+"22", sshConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	// create new SFTP client
	client, err := sftp.NewClient(connection)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// create destination file
	srcFile, err := client.Open(srcFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	// create source file
	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	// copy source file to destination file
	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("File %s: %d bytes copied\n", srcFilePath, bytes)
}

func getHostKey(host string) ssh.PublicKey {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		log.Fatalf("no hostkey found for %s", host)
	}

	return hostKey
}
