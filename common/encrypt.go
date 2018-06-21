package common

import "crypto/cipher"

//Cipher 加密
type Cipher struct {
	enc cipher.Stream
	dec cipher.Stream
}
