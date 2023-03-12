package gapi

import (
	db "IMChat/db/sqlc"
	"IMChat/pb"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 1.身份校验
// 2.将传进来的参数与审核表中的数据进行校验
// 3.校验成功后，往关系表中添加数据，并删除审核表中对应的数据
func (server *Server) AddContact(ctx context.Context, req *pb.AddContactRequest) (*pb.AddContactResponse, error) {
	_, err := server.authorization(ctx)
	if err != nil {
		return nil, err
	}

	if req.Type != 1 {
		return nil, status.Errorf(codes.Internal, "faild add friend")
	}

	// 获取审核表中的数据
	// select * from examine where owner_id = req.ownerId AND target_id = req.TargetId AND type = 1
	// if payload.Username != pb.GetUsername() {  }
	// 判断是否已是好友
	// 插入数据,并删除审核表中对应的数据
	_, err = server.store.AddContactTx(ctx, db.AddContactTxParams{
		AddContactParams: db.AddContactParams{
			OwnerID:  req.GetOwnerId(),
			TargetID: req.GetTargetId(),
			Type:     int16(req.GetType()),
		},
		AfterCreate: func(contact db.Contact) error {
			_, err := server.store.AddContactTx(context.Background(), db.AddContactTxParams{
				AddContactParams: db.AddContactParams{
					OwnerID:  req.GetTargetId(),
					TargetID: req.GetOwnerId(),
					Type:     int16(req.GetType()),
				},
				AfterCreate: func(contact db.Contact) error {
					return server.store.DeleteExamine(context.Background(), db.DeleteExamineParams{
						TargetID: contact.TargetID,
						OwnerID:  contact.OwnerID,
						Type:     contact.Type,
					})
				},
			})
			return err
		},
	})
	if err != nil {
		return nil, err
	}

	return &pb.AddContactResponse{Message: "Add Friend Success..."}, nil
}
