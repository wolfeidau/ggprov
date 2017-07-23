package ggprov

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// SSHSession state for the ssh session
type SSHSession struct {
	Client *ssh.Client
}

// NewSSHSession create an ssh session
func NewSSHSession(hostname, port, username, keyPath string) (*SSHSession, error) {

	log.Println("Loading key from path", keyPath)
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read key file")
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse private key")
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	log.Println("Connecting to host", hostname)

	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port), config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect")
	}

	return &SSHSession{Client: client}, nil
}

// Copy copy a file with the supplied name, mode and contents
func (ss *SSHSession) Copy(size int64, mode os.FileMode, fileName string, contents io.Reader, destinationPath string) error {
	return ss.copy(size, mode, fileName, contents, destinationPath)
}

// CopyPath copy a file
func (ss *SSHSession) CopyPath(filePath, destinationPath string) error {

	log.Println("CopyPath to host", filePath, destinationPath)

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return err
	}
	return ss.copy(s.Size(), s.Mode().Perm(), path.Base(filePath), f, destinationPath)
}

func (ss *SSHSession) copy(size int64, mode os.FileMode, fileName string, contents io.Reader, destination string) error {

	session, err := ss.Client.NewSession()
	if err != nil {
		return errors.Wrap(err, "Failed to create session for copy")
	}

	defer session.Close()
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		fmt.Fprintf(w, "C%#o %d %s\n", mode, size, fileName)
		io.Copy(w, contents)
		fmt.Fprint(w, "\x00")
	}()
	cmd := fmt.Sprintf("scp -t %s", destination)
	if err := session.Run(cmd); err != nil {
		return err
	}
	return nil
}
