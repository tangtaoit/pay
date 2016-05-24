package keys

import (
	"os"
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	"pay/config"
	"bufio"
)

func GetAlipayPrivateKey() *rsa.PrivateKey {

	privateKeyFile, err := os.Open(config.GetAlipaySetting().AlipayPrivateKeyPath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}
