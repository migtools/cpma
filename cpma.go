package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type Config struct {
	Source struct {
		User    string `json:"user"`
		CertKey string `json:"certKey"`
		Host    string `json:"host"`
	} `json:source`
}

func main() {
	config := LoadConfiguration("config.json")

	// get host public key
	hostKey := getHostKey(config.Source.Host)

	key, err := ioutil.ReadFile(config.Source.CertKey)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: config.Source.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},

		// allow any host key to be used (non-prod)
		//HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		// verify host public key
		HostKeyCallback: ssh.FixedHostKey(hostKey),

		Timeout: 10 * time.Second,
	}

	connection, err := ssh.Dial("tcp", config.Source.Host+":"+"22", sshConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	session, err := connection.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// setup standard out and error
	// uses writer interface
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Run("cat ./essay/config.yaml"); err != nil {
		panic("Failed to run: " + err.Error())
	}
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

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
