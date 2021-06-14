package account

import (
	"testing"
)


func TestPublicKeyAndAddress(t *testing.T) {
	publicKey := MustDerivePublicKey("c3af5ee3619783355c1504814e12230bdfc7b5e336f18fef336615e54881f4db")
	first := publicKey == "171363aaa6ddbfb5b3c1e6b24e8722518760f387a7133a8b472418f1fd320e96f62b4669753ae808af93653187f93fd4e972c1865ef023184408085cc40e007c"
	if !first {
		t.Logf("The public key is %s", publicKey)
		t.Error("publickey does not match")
	}

	compressedPublicKey := MustCompressPublicKey(publicKey)
	second := compressedPublicKey == "02171363aaa6ddbfb5b3c1e6b24e8722518760f387a7133a8b472418f1fd320e96"
	if !second {
		t.Logf("The compressed key is %s", compressedPublicKey)
		t.Error("Compressed Public Key does not match")
	}

	address := MustDeriveAddress(publicKey)
	third := address == "5dB994389C7495996A313F830f21015B9F8127B0"
	if !third {
		t.Logf("The address is %s", address)
		t.Error("Address does not match")
	}
}