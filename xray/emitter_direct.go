package xray

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"net"
)

type DirectEmitter struct {
	segmentProcessor *SegmentProcessor
}

func NewDirectEmitter(awsConfig aws.Config) *DirectEmitter {
	return &DirectEmitter{segmentProcessor:NewSegmentProcessor(awsConfig, 50)}
}

func (de *DirectEmitter) Emit(seg *Segment) {
	ctx := context.Background()
	batch := segmentPacking(seg)
	de.segmentProcessor.Transmit(&ctx, batch)
}

func (de *DirectEmitter) RefreshEmitterWithAddress(raddr *net.UDPAddr) {
	return
}

func segmentPacking(seg *Segment) (batch []string) {
	return
}
