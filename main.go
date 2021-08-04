package main

import (
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := xray.Configure(xray.Config{
		DaemonAddr:     "xray-daemon:2000",
		ServiceVersion: "1.0.0",
	}); err != nil {
		panic(err)
	}

	e := echo.New()
	hello := echo.WrapHandler(xray.Handler(xray.NewFixedSegmentNamer("test-app"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// http request
		req, err := http.NewRequest(http.MethodGet, "https://aws.amazon.com/", nil)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		resp, err := xray.Client(http.DefaultClient).Do(req.WithContext(ctx))
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		// sub segment
		subCtx, subSeg := xray.BeginSubsegment(ctx, "waiting-something")

		// aws sdk
		cfg, err := config.LoadDefaultConfig(subCtx)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		awsv2.AWSV2Instrumentor(&cfg.APIOptions)
		svc := s3.NewFromConfig(cfg)
		if _, err := svc.ListBuckets(subCtx, &s3.ListBucketsInput{}); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		subSeg.Close(nil)

		w.Write([]byte("Hello!"))
	})))
	e.GET("/", hello)
	e.Logger.Fatal(e.Start(":8080"))
}
