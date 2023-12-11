package v1

import (
	"context"
	"time"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (userGrpcServer *UserGrpcServer) VerifyEmail(ctx context.Context, request *pb.VerifyEmailRequest) (*pb.GenericResponse, error) {
	code := request.GetVerificationCode()

	verificationCode := commonUtility.Encode(code)

	query := bson.D{{Key: "verificationCode", Value: verificationCode}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: true}, {Key: "updated_at", Value: time.Now()}}}, {Key: "$unset", Value: bson.D{{Key: "verificationCode", Value: constants.EmptyString}}}}
	result, err := userGrpcServer.userCollection.UpdateOne(ctx, query, update)

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
