package sftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
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
	sshClient  *ssh.Client  `json:"ssh_client"`
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

	if params.sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}

	// create sftp client
	if params.sftpClient, err = sftp.NewClient(params.sshClient); err != nil {
		return err
	}
	return nil
}

func (params *Ssh) createRemoteDir(remoteDir string) error {
	session, err := params.sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	err = session.Run(fmt.Sprintf("mkdir -p %s", remoteDir))
	return err
}

func (params *Ssh) UploadFile(localFilePath, remoteDir string) error {
	err := params.createRemoteDir(remoteDir)
	if err != nil {
		return err
	}
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
		n, err := srcFile.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				break
			}
		}
		dstFile.Write(buf[:n])
	}
	return err
}

func (params *Ssh) Download(remotePath, localPath string) error {
	remoteFile, _ := params.sftpClient.Open(remotePath) //远程
	localFile, _ := os.Create(localPath)                //本地
	defer func() {
		_ = remoteFile.Close()
		_ = localFile.Close()
	}()

	if _, err := remoteFile.WriteTo(localFile); err != nil {
		return err
	}
	return nil
}

func (params *Ssh) Close() {
	params.sftpClient.Close()
	params.sshClient.Close()
}
