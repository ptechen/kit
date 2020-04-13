package sftp

import "testing"

func TestSsh_SendFile(t *testing.T) {
	ssh := &Ssh{
		User:     "root",
		Password: "123",
		Host:     "192.168.3.245",
		Port:     22,
	}
	defer ssh.Close()
	err := ssh.Connect()
	if err != nil {
		t.Error(err)
	}
	err = ssh.UploadFile("/Users/taochen/go/src/kit/sftp/sftp.go", "/datas/nfs")
	if err != nil {
		t.Error(err)
	}

	err = ssh.createRemoteDir("/data/data")
	if err != nil {
		t.Error()
	}
	err = ssh.checkFileExist("/data/data")
	if err != nil {
		t.Error()
	}

}
