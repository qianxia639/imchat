package gapi

import (
	"IMChat/pb"
	"io"
)

func (server *Server) SendMessage(stream pb.Message_SendMessageServer) error {
	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			in, err := stream.Recv()
			if err == io.EOF {
				return nil
			}

			if err != nil {
				return err
			}

			if err := stream.Send(&pb.SendMessageResponse{Msg: in.Msg}); err != nil {
				return err
			}
		}
	}

}
