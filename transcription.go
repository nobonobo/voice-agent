package main

import (
	"context"
	"io"
	"log"
	"sync"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "cloud.google.com/go/speech/apiv1/speechpb"
	"google.golang.org/api/option"
	"google.golang.org/grpc/status"
)

var gcpClient *speech.Client

func Transcription(ctx context.Context, fragments <-chan []byte, response chan<- *speechpb.StreamingRecognitionResult) error {
	if gcpClient == nil {
		c, err := speech.NewClient(ctx, option.WithAPIKey(config.GcpAPIKey))
		if err != nil {
			return err
		}
		gcpClient = c
	}
	stream, err := gcpClient.StreamingRecognize(ctx)
	if err != nil {
		return err
	}
	if err := stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					LanguageCode:    "ja-JP", //"ja-JP" "en-US"
					Encoding:        speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: 16000,
				},
				InterimResults: true,
			},
		},
	}); err != nil {
		return err
	}
	cancel := sync.OnceFunc(func() { stream.CloseSend() })
	defer cancel()
	go func() {
		defer cancel()
		for v := range fragments {
			req := &speechpb.StreamingRecognizeRequest{
				StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
					AudioContent: v,
				},
			}
			if err := stream.Send(req); err != nil {
				if s, ok := status.FromError(err); ok && s.Code() == 13 {
					return
				}
				log.Println(err)
				return
			}
		}
	}()
	for {
		v, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if v.Error != nil {
			return err
		}
		for _, r := range v.Results {
			if r.GetStability() > float32(config.SttStabirity) {
				response <- r
				return nil
			}
		}
	}
	return nil
}
