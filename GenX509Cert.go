/*
X.509证书
PKI是公匙基础设施，是一个公匙集框架，这些公匙包含所有者的名称和位置等附加信息，以及之间链接能提供某种审批机制
目前主流的PKI是基于X.509证书。例如web浏览器使用它们来验证web站点的身份。
如下代码是为自己网站生产一个自签名x.509证书并保存在.cer文件中的例子。
*/
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

func main() {
	random := rand.Reader

	var key rsa.PrivateKey
	loadKey("./ca/private.key", &key)

	now := time.Now()
	then := now.Add(60 * 60 * 24 * 365 * 1000 * 1000 * 1000) //一年有效期
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "jun.ming.wang",
			Organization: []string{"zhejianglab"},
		},
		NotBefore:    now,
		NotAfter:     then,
		SubjectKeyId: []byte{1, 2, 3, 4},
		KeyUsage:     x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,

		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              []string{"jun.ming.wang", "localhost"},
	}
	derBytes, err := x509.CreateCertificate(random, &template, &template, &key.PublicKey, &key)
	checkError(err)

	certFile, err := os.Create("./ca/jun.ming.wang.cer")
	checkError(err)
	certFile.Write(derBytes)
	certFile.Close()

	certPEMFile, err := os.Create("./ca/jun.ming.wang.pem")
	checkError(err)
	pem.Encode(certPEMFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certPEMFile.Close()

	keyPEMFile, err := os.Create("./ca/private.pem")
	checkError(err)
	pem.Encode(keyPEMFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(&key)})
	keyPEMFile.Close()
}

func loadKey(fileName string, key interface{}) {
	inFile, err := os.Open(fileName)
	checkError(err)
	decoder := gob.NewDecoder(inFile)
	err = decoder.Decode(key)
	checkError(err)
	inFile.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
