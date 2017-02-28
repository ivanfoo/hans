package utils

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var netrcFile = filepath.Join(os.Getenv("HOME"), ".netrc")

type netrcCredentials struct {
	machine  string
	login    string
	password string
}

func SetGitCredentials(gitRemote string, base64Credentials string) error {
	var n netrcCredentials
	n.parseServerName(gitRemote)
	n.decodeCredentials(base64Credentials)
	n.saveToFile()
	return nil
}

func (n *netrcCredentials) parseServerName(gitRemote string) {
	u, err := url.Parse(gitRemote)
	if err != nil {
		log.Fatal(err)
	}
	n.machine = u.Host
}

func (n *netrcCredentials) decodeCredentials(base64Credentials string) {
	data, err := base64.StdEncoding.DecodeString(base64Credentials)

	if err != nil {
		log.Fatal(err)
	}

	credentials := strings.Split(string(data), ":")
	n.login = credentials[0]
	n.password = credentials[1]
}

func (n *netrcCredentials) saveToFile() error {
	netrcConfig := fmt.Sprintf("machine %s\nlogin %s\npassword %s", n.machine, n.login, n.password)
	err := ioutil.WriteFile(netrcFile, []byte(netrcConfig), 0600)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
