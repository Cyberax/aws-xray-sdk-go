package xray

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/xray"
	"sync"
)

// transmit a single batch to XRay
func EmitSegments(client *xray.Client, ctx *context.Context, batch *[]string, group sync.WaitGroup) {
	req := &xray.PutTraceSegmentsInput{TraceSegmentDocuments: *batch}
	putSegReq := client.PutTraceSegmentsRequest(req)
	_, _ = putSegReq.Send(*ctx)

	group.Done()
}
