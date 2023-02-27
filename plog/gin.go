package plog

//
//// gin 框架中使用log
//
//import (
//	"git.im30.lan/kop/common/trace"
//	"github.com/gin-gonic/gin"
//	"time"
//)
//
//const (
//	ContextLogKey = "__context_log_key__"
//)
//
//const (
//	requestStartTimeKey = "start"
//	traceIDKey          = "trace"
//	queryPathKey        = "path"
//	nameKey             = "name"
//)
//
//func GetDefaultFieldEntry(ctx *gin.Context) *Entry {
//	logEntry := getFromGin(ctx)
//	if logEntry != nil {
//		return logEntry
//	}
//
//	// 初始化
//	traceID := ""
//	if ctx != nil && ctx.Request != nil && ctx.Request.Header != nil {
//		traceID = ctx.GetHeader(trace.HeaderTraceIdKey)
//		// 正常情况接收到的 traceid 不会很长
//		if len(traceID) > 50 {
//			traceID = traceID[:50]
//		}
//	}
//	if traceID == "" {
//		traceID = trace.ID()
//	}
//	path := ""
//	if ctx != nil && ctx.Request != nil && ctx.Request.URL != nil {
//		path = ctx.Request.URL.Path
//	}
//
//	// 默认字段
//	fields := map[string]interface{}{
//		requestStartTimeKey: time.Now().Format(logTimeFormatter),
//		traceIDKey:          traceID,
//		queryPathKey:        path,
//	}
//	logEntry = stdLogger.withFields(fields)
//
//	// 保存上下文
//	if ctx != nil {
//		ctx.Set(ContextLogKey, logEntry)
//		ctx.Set(trace.ContextTraceId, traceID)
//	}
//	if ctx != nil && ctx.Request != nil && ctx.Request.Header != nil {
//		ctx.Header(trace.HeaderTraceIdKey, traceID)
//	}
//	return logEntry
//}
//
//// GetFromGin return a Entry, you should use Entry in api handler to log
//func getFromGin(c *gin.Context) *Entry {
//	var (
//		e  *Entry
//		ok bool
//	)
//	if c == nil {
//		return nil
//	}
//
//	ee, ok := c.Get(ContextLogKey)
//	if ok {
//		e, ok = ee.(*Entry)
//	}
//	if ok {
//		return e
//	}
//	return nil
//}
//
//// warp stdLogger
//func Debug(ctx *gin.Context, args ...interface{}) {
//	GetDefaultFieldEntry(ctx).Debug(args...)
//	stdLogger.logger.Debug(args...)
//}
//
//func Debugf(ctx *gin.Context, format string, args ...interface{}) {
//	GetDefaultFieldEntry(ctx).Debugf(format, args...)
//}
//
//func Info(ctx *gin.Context, args ...interface{}) {
//	GetDefaultFieldEntry(ctx).Info(args...)
//}
//
//func Infof(ctx *gin.Context, format string, args ...interface{}) {
//	GetDefaultFieldEntry(ctx).Infof(format, args...)
//}
//
//func Warn(ctx *gin.Context, args ...interface{}) {
//	GetDefaultFieldEntry(ctx).Warn(args...)
//}
//
//func Warnf(ctx *gin.Context, format string, args ...interface{}) {
//	GetDefaultFieldEntry(ctx).Warnf(format, args...)
//}
//
//func Error(ctx *gin.Context, args ...interface{}) {
//	GetDefaultFieldEntry(ctx).Error(args...)
//}
//
//func Errorf(ctx *gin.Context, format string, args ...interface{}) {
//	GetDefaultFieldEntry(ctx).Errorf(format, args...)
//}
//
//// Fatal will call os.Exit(1), be careful to use.
//func Fatal(ctx *gin.Context, args ...interface{}) {
//	GetDefaultFieldEntry(ctx).Fatal(args...)
//}
//
//// Fatalf will call os.Exit(1), be careful to use.
//func Fatalf(ctx *gin.Context, format string, args ...interface{}) {
//	GetDefaultFieldEntry(ctx).Fatalf(format, args...)
//}
