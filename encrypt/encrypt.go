package encrypt

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash/crc32"
	"github.com/OneOfOne/xxhash"
)

// 生成md5
func MD5(str string) string {
	c := md5.New()
	c.Write([]byte(str))
	return hex.EncodeToString(c.Sum(nil))
}

//生成sha1
func SHA1(str string) string {
	c := sha1.New()
	c.Write([]byte(str))
	return hex.EncodeToString(c.Sum(nil))
}

func CRC32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

// RS Hash Function
func RSHash(str []byte) uint32 {
	var b uint32 = 378551
	var a uint32 = 63689
	var hash uint32 = 0
	for i := 0; i < len(str); i++ {
		hash = hash*a + uint32(str[i])
		a *= b
	}
	return hash
}

func Hash(input []byte) uint64 {
	h := xxhash.New64()
	h.Write(input)
	return h.Sum64()
}
