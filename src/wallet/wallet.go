package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"utils"
	"golang.org/x/crypto/ripemd160"
	"crypto/sha256"
)

const version = byte(0x00)
const walletFile = "wallet.dat"
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallet() *Wallet {
	privateKey, publicKey := newKeyPair()
	return &Wallet{privateKey, publicKey}
}

/**
	生成一对公、私钥
 */
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	utils.LogE(err)
	return *private, pubKey
}

/**
	公钥生成地址
 */
func (wallet Wallet) GetAddress() []byte {
	pubKeyHash := hashPubKey(wallet.PublicKey)
	checkSum := checkSum(pubKeyHash)
	fullAddressData := append([]byte{version}, pubKeyHash...)
	fullAddressData = append(fullAddressData, checkSum...)
	return utils.Base58Encode(fullAddressData)
}


func checkSum(pubKeyHash []byte) []byte {
	doubleHash := sha256.Sum256(sha256.Sum256(pubKeyHash)[:])
	return doubleHash[:4]
}

/**
	对公钥两次哈希
 */
func hashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	utils.LogE(err)
	return publicRIPEMD160
}
