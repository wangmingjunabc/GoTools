/*
读取X.509证书
*/
package main

import (
	"crypto/x509"
	"fmt"
	"os"
)

func main() {
	certCerFile, err := os.Open("./ca/jun.ming.wang.cer")
	checkError(err)
	derBytes := make([]byte, 1000) //比证书文件大
	count, err := certCerFile.Read(derBytes)
	checkError(err)
	certCerFile.Close()

	//截取到证书实际长度
	cert, err := x509.ParseCertificate(derBytes[0:count])
	checkError(err)

	fmt.Println("Name %s\n", cert.Subject.CommonName)
	fmt.Println("Not before %s\n", cert.NotBefore.String())
	fmt.Println("Not after %s\n", cert.NotAfter.String())
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
/*
output:
Name jun.ming.wang
Not before 2021-02-03 06:47:17 +0000 UTC
Not after 2022-02-03 06:47:17 +0000 UTC
*/
