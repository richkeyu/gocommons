package wrapper

import (
	"context"

	"github.com/richkeyu/gocommons/client"
	"github.com/richkeyu/gocommons/server"
	"github.com/richkeyu/gocommons/trace"
)

func HttpClientTrace(next client.Wrapper) client.Wrapper {
	return func(ctx context.Context, req *client.Request) (*client.Response, error) {
		// 获取traceId
		traceID := ""
		request := server.FromContext(ctx)
		if request != nil {
			traceID = request.Header(trace.HeaderTraceIdKey)
			// 正常情况接收到的 traceid 不会很长
			if len(traceID) > 50 {
				traceID = traceID[:50]
			}
		}
		if traceID == "" {
			traceID = trace.ID()
			// 保存回去 供后续使用
			ginCtx := server.GinFromContext(ctx)
			if ginCtx != nil {
				// 此处依赖了gin.Context
				if ginCtx != nil && ginCtx.Request != nil && ginCtx.Request.Header != nil {
					ginCtx.Header(trace.HeaderTraceIdKey, traceID)
				}
			}
		}
		// 注入请求
		req.GetRequest().Header.Add(trace.HeaderTraceIdKey, traceID)

		return next(ctx, req)
	}
}
