package middleware

import (
	"bytes"
	"donate/logger"
	"io"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func GinLogger(Lg *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		bodyBuf := bufferPool.Get().(*bytes.Buffer)
		defer func() {
			bodyBuf.Reset()
			bufferPool.Put(bodyBuf)
		}()

		if c.Request.Body != nil {
			_, err := io.Copy(bodyBuf, c.Request.Body)
			if err != nil {
				log.Error().Err(err).Msg("Failed to read request body")
				c.Error(err)
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBuf.Bytes()))
		}

		defer func() {
			ctxLogger := c.MustGet(DefaultLoggerKey).(*logger.CtxLogger)
			now := time.Since(start).Milliseconds()

			// 如果log level > info, 则记录 request_body
			logLevel := ctxLogger.GetLevel()
			if logLevel > zerolog.InfoLevel {
				bodyBuf.Reset()
				bufferPool.Put(bodyBuf)
			}

			logEvent := ctxLogger.Info().
				Int("status", c.Writer.Status()).
				Str("method", c.Request.Method).
				Str("path", path).
				Str("query", query).
				Str("ip", c.ClientIP()).
				Str("user-agent", c.Request.UserAgent()).
				Str("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()).
				Int64("cost(ms)", now)

			if bodyBuf.Len() > 0 {
				logEvent = logEvent.RawJSON("request_body", bodyBuf.Bytes())
			}

			// 如果log level > info, 则发送,否则不发送
			if logLevel > zerolog.InfoLevel {
				logEvent.Send()
			}
		}()

		c.Next()
	}
}
