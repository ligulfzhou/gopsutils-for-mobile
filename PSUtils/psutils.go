package PSUtils

import (
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	CipherList = []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"}
)

type PSUtils struct {
	user, password, host, keyPath, keyString string
	port                                     int
	cipherList                               []string

	platform string

	// for network
	NetworkInterface                 string
	RX_LAST_TMSTAMP, TX_LAST_TMSTAMP int64
	RX_LAST_TOTAL, TX_LAST_TOTAL     int64

	// disk
	STORAGE_DEVICE_NAMES  []string
	PROC_DISKSTAT_TMSTAMP int64
	LastDiskStat          ProcDiskStats
	// LastReadsCompletedSuccess       int64
	// LastReadsMerged                 int64
	// LastSectorsRead                 int64
	// LastTimeSpentReadingMS          int64
	// LastWritesCompleted             int64
	// LastWritesMerged                int64
	// LastSectorsWritten              int64
	// LastTimeSpentWritingMS          int64
	// LastIOsCurrentlyInProgress      int64
	// LastTimeSpentDoingIOsMS         int64
	// LastWeightedTimeSpentDoingIOsMS int64

	// kernel version, useful for calc disk stat
	KERNEL_VERSION string

	client *ssh.Client
}

//func NewPSUtils(user, password, host, key string, port int, cipherList []string) *PSUtils {
func NewPSUtils(user, password, host, keyPath, keyString string, port int) *PSUtils {

	return &PSUtils{
		user:       user,
		password:   password,
		host:       host,
		keyPath:    keyPath,
		keyString:  keyString,
		port:       port,
		cipherList: CipherList,
		client:     nil,
	}
}

func (ps *PSUtils) Connect() (bool, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	if ps.keyPath == "" && ps.keyString == "" {
		auth = append(auth, ssh.Password(ps.password))
	} else {
		var pemBytes []byte
		if ps.keyPath != "" {
			pemBytes, err = ioutil.ReadFile(ps.keyPath)
			if err != nil {
				return false, err
			}
		} else {
			pemBytes = []byte(ps.keyString)
		}

		var signer ssh.Signer
		if ps.password == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(ps.password))
		}
		if err != nil {
			return false, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	config = ssh.Config{
		Ciphers: ps.cipherList,
	}

	clientConfig = &ssh.ClientConfig{
		User: ps.user,
		Auth: auth,
		// 5 second shall be acceptable.
		Timeout:         5 * time.Second,
		Config:          config,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr = fmt.Sprintf("%s:%d", ps.host, ps.port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return false, err
	}
	ps.client = client
	return true, nil
}

func (ps *PSUtils) checkConn() {
	select {}
}
