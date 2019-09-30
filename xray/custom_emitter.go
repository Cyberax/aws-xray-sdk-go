package xray

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/xray"
	log "github.com/sirupsen/logrus"
	"sync"
)

// transmit a single batch to XRay
func EmitSegments(client *xray.Client, ctx *context.Context, batch *[]string, group sync.WaitGroup) {
	req := &xray.PutTraceSegmentsInput{TraceSegmentDocuments: *batch}
	putSegReq := client.PutTraceSegmentsRequest(req)
	_, err := putSegReq.Send(*ctx)
	if err != nil {
		log.Errorf(err.Error())
	}

	group.Done()
}
