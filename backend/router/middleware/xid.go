package middleware

import (
	"encoding/base64"
	"encoding/binary"
	"os"
	"time"

	"donate/clock"
	iLog "donate/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var (
	clk    clock.Clock = clock.New()
	loc, _             = time.LoadLocation("Asia/Shanghai")
)

func GinXid(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		xid := GenReqId()

		log := *logger
		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str(DefaultXid, xid)
		})

		Log := iLog.NewCtxLogger(&log, iLog.APISimulation)
		c.Header(DefaultXid, xid)
		c.Set(DefaultLoggerKey, Log)
		c.Set(DefaultXid, xid)

		c.Next()
	}
}

var (
	pid = uint32(os.Getpid())
)

func GenReqId() string {
	prefix := clk.Now().In(loc).Format("20060102150405")
	var b [12]byte
	binary.LittleEndian.PutUint32(b[:], pid)
	binary.LittleEndian.PutUint64(b[4:], uint64(time.Now().UnixNano()))
	return prefix + base64.URLEncoding.EncodeToString(b[:])
}

func GetOrGenXid(c *gin.Context) string {
	xid, ok := c.Get(DefaultXid)
	if !ok {
		xid = GenReqId()
		c.Set(DefaultXid, xid)
	}
	return xid.(string)
}
