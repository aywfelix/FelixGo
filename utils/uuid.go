package utils

import (
	"encoding/binary"
	"encoding/hex"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/felix/felixgo/encrypt"
)

var (
	globalID uint64 = 0
)

type UniqueCode struct {
	endian    binary.ByteOrder
	hostID    uint32
	processId uint32
	baseID    uint32
}

func NewUniqueCode() *UniqueCode {
	uniq := &UniqueCode{}
	uniq.endian = binary.LittleEndian
	uniq.hostID = encrypt.RSHash([]byte(GetLocalIPV4()))
	uniq.processId = uint32(os.Getpid())
	uniq.baseID = 0
	return uniq
}

func (u *UniqueCode) NewUniqID() string {
	objID := make([]byte, 16)
	u.endian.PutUint32(objID[0:], uint32(time.Now().Unix()))
	u.endian.PutUint32(objID[4:], u.hostID)
	u.endian.PutUint32(objID[8:], u.processId)
	u.endian.PutUint32(objID[12:], atomic.AddUint32(&u.baseID, 1))
	return hex.EncodeToString(objID)
}

func GetGUID() uint64 {
	return atomic.AddUint64(&globalID, 1)
}

func UniqID() string {
	return strconv.FormatInt(int64(SnowFlake.GenInt()), 10)
}

// 基于服务器生成唯一ID
func GenUID(serverID uint64) uint64 {
	guid := GetGUID() % 1000 //此处需要处理 15毫秒内同时生成1000个的情况
	return serverID*1000000000000000 + uint64(TimeMillisecond()*1000)%1000000000000*1000 + guid
}
