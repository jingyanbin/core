package xnet

import "encoding/binary"

func CountBCC(buf []byte, offset int, length int) byte {
	value := option.XORBcc
	for i := offset; i < offset+length; i++ {
		value ^= buf[i]
	}
	return value
}

func XOREncrypt(seed uint32, buf []byte, offset int, length int) {
	key := make([]byte, 4)
	binary.LittleEndian.PutUint32(key, seed)
	var k int
	var c, x byte
	for i := offset; i < offset+length; i++ {
		k &= 3 //= k %= 4
		x = (buf[i] ^ key[k]) + c
		k++
		c = x
		buf[i] = x
	}
	return
}

func XORDecrypt(seed uint32, buf []byte, offset int, length int) {
	key := make([]byte, 4)
	binary.LittleEndian.PutUint32(key, seed)
	var k int
	var c, x byte
	for i := offset; i < offset+length; i++ {
		k &= 3 //= k %= 4
		x = (buf[i] - c) ^ key[k]
		k++
		c = buf[i]
		buf[i] = x
	}
}

type XORCrypt struct {
	iSeed uint32
	oSeed uint32
}

func (m *XORCrypt) Encrypt(buf []byte, offset int, length int) {
	m.oSeed = m.oSeed*option.XORCryptA + option.XORCryptB
	XOREncrypt(m.oSeed, buf, offset, length)
}

func (m *XORCrypt) Decrypt(buf []byte, offset int, length int) {
	m.iSeed = m.iSeed*option.XORCryptA + option.XORCryptB
	XORDecrypt(m.iSeed, buf, offset, length)
}

func NewXORCrypt(iSeed, oSeed uint32) *XORCrypt {
	return &XORCrypt{iSeed: iSeed, oSeed: oSeed}
}
