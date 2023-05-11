package v1

import (
	"context"
	"time"

	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (userServer *UserServer) VerifyEmail(ctx context.Context, request *pb.VerifyEmailRequest) (*pb.GenericResponse, error) {
	code := request.GetVerificationCode()

	verificationCode := utility.Encode(code)

	query := bson.D{{Key: "verificationCode", Value: verificationCode}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: true}, {Key: "updated_at", Value: time.Now()}}}, {Key: "$unset", Value: bson.D{{Key: "verificationCode", Value: ""}}}}
	result, err := userServer.userCollection.UpdateOne(ctx, query, update)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if result.MatchedCount == 0 {
		return nil, status.Errorf(codes.PermissionDenied, "Could not verify email address")
	}

	res := &pb.GenericResponse{
		Status:  "success",
		Message: "Email verified successfully",
	}

	return res, nil
}
