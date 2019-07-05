package cryptoHelper

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

func cipher2CipherFunc(cipher string) packet.CipherFunction {
	switch cipher {
	default:
		if cipher != "" {
			fmt.Println("Invalid cipher, using default")
		}
		fallthrough
	case "aes256":
		return packet.CipherAES256
	case "aes192":
		return packet.CipherAES192
	case "aes128":
		return packet.CipherAES128
	}
}

func Encrypt(plaintext []byte, password []byte, packetConfig *packet.Config) (ciphertext []byte, err error) {
	encbuf := bytes.NewBuffer(nil)

	w, err := armor.Encode(encbuf, "PGP MESSAGE", nil)
	if err != nil {
		return
	}
	defer w.Close()

	pt, err := openpgp.SymmetricallyEncrypt(w, password, nil, packetConfig)
	if err != nil {
		return
	}
	defer pt.Close()

	_, err = pt.Write(plaintext)
	if err != nil {
		return
	}

	// Close writers to force-flush their buffer
	pt.Close()
	w.Close()
	ciphertext = encbuf.Bytes()

	return
}

func Decrypt(ciphertext []byte, password []byte, packetConfig *packet.Config) (plaintext []byte, err error) {
	decbuf := bytes.NewBuffer(ciphertext)

	armorBlock, err := armor.Decode(decbuf)
	if err != nil {
		return
	}

	failed := false
	prompt := func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		// If the given passphrase isn't correct, the function will be called again, forever.
		// This method will fail fast.
		// Ref: https://godoc.org/golang.org/x/crypto/openpgp#PromptFunction
		if failed {
			return nil, errors.New("decryption failed")
		}
		failed = true
		return password, nil
	}

	md, err := openpgp.ReadMessage(armorBlock.Body, nil, prompt, packetConfig)
	if err != nil {
		return
	}

	plaintext, err = ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return
	}

	return
}

func EncryptBytes(origBytes []byte, pubKeyringFile string) ([]byte, error) {
	keyringFileBuffer, _ := os.Open(pubKeyringFile)
	defer keyringFileBuffer.Close()
	entityList, err := openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		return nil, err
	}

	// encrypt string
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, entityList, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(origBytes)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	encryptedBytes, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	return encryptedBytes, nil
}

func DecryptBytesWithIoReader(encBytes []byte, keyringFileBuffer io.Reader, passphrase string) ([]byte, error) {
	// init some vars
	var entity *openpgp.Entity
	var entityList openpgp.EntityList

	// Open the private key file
	entityList, err := openpgp.ReadArmoredKeyRing(keyringFileBuffer)
	//entityList, err := openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		return nil, err
	}
	entity = entityList[0]

	// Get the passphrase and read the private key.
	// Have not touched the encrypted string yet
	passphraseByte := []byte(passphrase)

	err = entity.PrivateKey.Decrypt(passphraseByte)
	if err != nil {
		return nil, err
	}

	for _, subkey := range entity.Subkeys {
		err = subkey.PrivateKey.Decrypt(passphraseByte)
		if err != nil {
			return nil, err
		}
	}

	block, err := armor.Decode(bytes.NewBuffer(encBytes))
	// Decrypt it with the contents of the private key
	md, err := openpgp.ReadMessage(block.Body, entityList, nil, nil)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func DecryptBytes(encBytes []byte, privKeyringFile string, passphrase string) ([]byte, error) {
	// init some vars
	var entity *openpgp.Entity
	var entityList openpgp.EntityList

	// Open the private key file
	keyringFileBuffer, err := os.Open(privKeyringFile)
	if err != nil {
		return nil, err
	}
	defer keyringFileBuffer.Close()
	entityList, err = openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		return nil, err
	}
	entity = entityList[0]

	// Get the passphrase and read the private key.
	// Have not touched the encrypted string yet
	passphraseByte := []byte(passphrase)
	log.Println("Decrypting private key using passphrase")

	entity.PrivateKey.Decrypt(passphraseByte)
	for _, subkey := range entity.Subkeys {
		subkey.PrivateKey.Decrypt(passphraseByte)
	}

	log.Println("Finished decrypting private key using passphrase")

	// Decrypt it with the contents of the private key
	md, err := openpgp.ReadMessage(bytes.NewBuffer(encBytes), entityList, nil, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func DecryptBase64(enc string, privateKey string, passphrase string) ([]byte, error) {
	content, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return nil, err
	}

	privateKeyReader := strings.NewReader(privateKey)
	return DecryptBytesWithIoReader(content, privateKeyReader, passphrase)
}
