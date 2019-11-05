package xray

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/xray"
	"google.golang.org/appengine/log"
	"sync"
)

type SegmentProcessor struct {
	batchSize int
	XRayClient *xray.Client
}

func NewSegmentProcessor(awsConfig aws.Config, batchSize int) *SegmentProcessor {
	processor := SegmentProcessor{
		batchSize: batchSize,
	}

	processor.XRayClient = xray.New(awsConfig)
	return &processor
}

func (p *SegmentProcessor) transform(batch []string) [][]string {
	batches := make([][]string, 0)
	curBatch := make([]string, 0)
	for idx, s := range batch {
		if idx > 0 && idx % p.batchSize == 0 {
			batches = append(batches, curBatch)
			curBatch = make([]string, 0)
		}
		curBatch = append(curBatch, s)
	}
	return batches
}

func (p *SegmentProcessor) Transmit(ctx *context.Context, batch []string) {
	batches := p.transform(batch)
	numProcessors := len(batches)
	var group sync.WaitGroup
	group.Add(numProcessors)
	for i := 0; i < numProcessors; i++ {
		go EmitSegments(p.XRayClient, ctx, &(batches[i]), group)
	}
	group.Wait()
}

// transmit a single batch to XRay
func EmitSegments(client *xray.Client, ctx *context.Context, batch *[]string, group sync.WaitGroup) {
	req := &xray.PutTraceSegmentsInput{TraceSegmentDocuments: *batch}
	putSegReq := client.PutTraceSegmentsRequest(req)
	_, err := putSegReq.Send(*ctx)
	if err != nil {
		log.Errorf(*ctx, err.Error())
	}

	group.Done()
}
