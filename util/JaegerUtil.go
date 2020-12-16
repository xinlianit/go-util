package util
//
//import (
//	"context"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/micro/go-micro/v2/logger"
//	"github.com/micro/go-micro/v2/metadata"
//	"github.com/opentracing/opentracing-go"
//	"github.com/opentracing/opentracing-go/log"
//	"github.com/uber/jaeger-client-go"
//	jaegerConfig "github.com/uber/jaeger-client-go/config"
//	"io"
//)
//
//var JaegerTracer = struct {
//	KeyRootSpan     string
//	KeyParentSpan   string
//}{
//	KeyRootSpan:     "RootSpan",
//	KeyParentSpan:   "ParentSpan",
//}
//
//type JaegerTracerConfig struct {
//	LocalAgentHostPort string
//	CollectorEndpoint  string
//}
//
//type JaegerUtil struct {
//	ctx *gin.Context
//	cfg jaegerConfig.Configuration
//}
//
//// 创建追踪器
//func NewJaegerTracer(ctx *gin.Context) *JaegerUtil {
//	return &JaegerUtil{ctx: ctx}
//}
//
//// 配置追踪器
//func (u *JaegerUtil) Config(config JaegerTracerConfig) {
//	// 初始化配置
//	u.cfg = jaegerConfig.Configuration{
//		//ServiceName: "", // 服务名
//		//Disabled:    false,
//		//RPCMetrics:  false,
//		//Tags:        nil,
//
//		// 采样配置
//		Sampler: &jaegerConfig.SamplerConfig{
//			// todo 采样类型说明文档
//			Type:  jaeger.SamplerTypeConst, // 采样类型：const-固定采样
//			Param: 1,                       // 采样参数：1-全采样、0-不采样
//			//SamplingServerURL:        "",
//			//SamplingRefreshInterval:  0,
//			//MaxOperations:            0,
//			//OperationNameLateBinding: false,
//			//Options:                  nil,
//		},
//
//		// 上报配置
//		Reporter: &jaegerConfig.ReporterConfig{
//			//QueueSize:                  0,
//			//BufferFlushInterval:        0,
//			LogSpans: true,
//			//LocalAgentHostPort:         config.LocalAgentHostPort, // Jaeget Agent 代理地址  todo LocalAgentHostPort 是否可以使用远程代理
//			//DisableAttemptReconnecting: false,
//			//AttemptReconnectInterval:   0,
//			CollectorEndpoint: config.CollectorEndpoint, // 将span发往jaeger-collector的服务地址 todo LocalAgentHostPort && CollectorEndpoint 区别
//			//User:                       "",
//			//Password:                   "",
//			//HTTPHeaders:                nil,
//		},
//
//
//		//Headers:             nil,
//		//BaggageRestrictions: nil,
//		//Throttler:           nil,
//	}
//}
//
//// 创建全局追踪器
//func (u *JaegerUtil) NewGlobalTracer(serviceName string) (opentracing.Tracer, io.Closer) {
//
//	// 设置服务名
//	u.cfg.ServiceName = serviceName
//
//	// 创建跟踪器
//	tracer, closer, err := u.cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
//
//	if err != nil {
//		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
//	}
//
//	// 注册跟踪器
//	opentracing.SetGlobalTracer(tracer)
//
//	return tracer, closer
//}
//
//// 注入追踪跨度
//func JaegerTracerInject(span opentracing.Span) context.Context {
//	// 注入
//	md := make(map[string]string)
//	errInject := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
//	if errInject != nil {
//		logger.Errorf("[JaegerMiddleware] Inject Error: %v", errInject.Error())
//		span.SetTag("error", true)
//		span.LogFields(log.String("inject-error", errInject.Error()))
//	}
//
//	ctx := context.TODO()
//	ctx = opentracing.ContextWithSpan(ctx, span)
//	ctx = metadata.NewContext(ctx, md) // 设置数据到上下文
//
//	return ctx
//}
//
//// 设置根跨度
//func (u *JaegerUtil) SetRootSpan(span opentracing.Span) {
//	u.ctx.Set(JaegerTracer.KeyRootSpan, span)
//}
//
//// 获取根跨度
//func (u *JaegerUtil) GetRootSpan() (span opentracing.Span, ok bool) {
//	val, exists := u.ctx.Get(JaegerTracer.KeyRootSpan)
//	if exists == false {
//		ok = false
//		return
//	}
//
//	span, ok = val.(opentracing.Span)
//	return
//}
//
//// 设置父跨度
//func (u *JaegerUtil) SetParentSpan(span opentracing.Span) {
//	u.ctx.Set(JaegerTracer.KeyParentSpan, span)
//}
//
//// 获取父跨度
//func (u *JaegerUtil) GetParentSpan() (span opentracing.Span, ok bool) {
//	val, exists := u.ctx.Get(JaegerTracer.KeyParentSpan)
//	if exists == false {
//		ok = false
//		return
//	}
//
//	span, ok = val.(opentracing.Span)
//	return
//}
