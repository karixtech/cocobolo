package mock_cocobolo_test

import (
	"context"
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	cbpb "github.com/tsudot/cocobolo/cocobolo"
	cbmock "github.com/tsudot/cocobolo/mock_cocobolo"
)

var msg = &cbpb.CallbackRequest{
	URL:       "http://tsudot.com/",
	Method:    "GET",
	RequestId: "1",
}

var callbackResponse = &cbpb.CallbackResponse{
	RequestId: "1",
}

func TestCocobolo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stream := cbmock.NewMockCocobolo_MakeRequestClient(ctrl)

	stream.EXPECT().Send(msg).Return(nil)

	stream.EXPECT().Recv().Return(callbackResponse, nil)

	stream.EXPECT().CloseSend().Return(nil)

	cbclient := cbmock.NewMockCocoboloClient(ctrl)

	cbclient.EXPECT().MakeRequest(
		gomock.Any(),
	).Return(stream, nil)

	if err := testCocobolo(cbclient); err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

func testCocobolo(client cbpb.CocoboloClient) error {
	stream, err := client.MakeRequest(context.Background())

	if err != nil {
		return err
	}

	if err := stream.Send(msg); err != nil {
		return err
	}

	if err := stream.CloseSend(); err != nil {
		return err
	}

	got, err := stream.Recv()

	log.Printf("%v", got)

	return nil
}
