package sftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"path"
	"time"
)

type Ssh struct {
	User       string       `json:"user"`
	Password   string       `json:"password"`
	Host       string       `json:"host"`
	Port       int          `json:"port"`
	sftpClient *sftp.Client `json:"sftp_client"`
}

// UserInterface is test interface.
type SshInterface interface {
	Connect() error
	SendFile(localFilePath, remoteDir string) error
}

func (params *Ssh) Connect() error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(params.Password))

	clientConfig = &ssh.ClientConfig{
		User:    params.User,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connect to ssh
	addr = fmt.Sprintf("%s:%d", params.Host, params.Port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return err
	}
	params.sftpClient = sftpClient
	return nil
}

func (params *Ssh) SendFile(localFilePath, remoteDir string) error {
	defer params.sftpClient.Close()
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(localFilePath)
	dstFile, err := params.sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}
	return err
}
