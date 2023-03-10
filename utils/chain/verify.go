package chain

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func Verify(address string, salt string, signature string, wallet string) (bool, error) {
	var walletToStr = map[string]string{
		"metamask": "\x19Ethereum Signed Message:\n",
		"kaikas":   "\x19Klaytn Signed Message:\n",
	}
	hashed := []byte(walletToStr["wallet"] + strconv.Itoa(len(salt)) + salt)
	hash := crypto.Keccak256Hash(hashed)

	decoded := hexutil.MustDecode(signature)
	if decoded[64] == 27 || decoded[64] == 28 {
		decoded[64] -= 27
	}

	sigPublickKeyECDSA, err := crypto.SigToPub(hash.Bytes(), decoded)

	if sigPublickKeyECDSA == nil {
		err = errors.New("no pub key found")
	}
	if err != nil {
		return false, err
	}
	actualAddr := crypto.PubkeyToAddress(*sigPublickKeyECDSA).String()
	if strings.EqualFold(actualAddr, address) {
		return true, nil
	} else {
		return false, errors.New("address doesn't match")
	}
}
