//nolint:nolintlint,dupl
package api

import (
	"context"
	"fmt"

	converter "github.com/NpoolPlatform/review-manager/pkg/converter/review"
	crud "github.com/NpoolPlatform/review-manager/pkg/crud/review"
	commontracer "github.com/NpoolPlatform/review-manager/pkg/tracer"
	tracer "github.com/NpoolPlatform/review-manager/pkg/tracer/review"

	constant "github.com/NpoolPlatform/review-manager/pkg/message/const"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/review/mgr/v2"

	"github.com/google/uuid"
)

func ValidateCreate(in *npool.ReviewReq) error {
	if in.ID != nil {
		if _, err := uuid.Parse(in.GetID()); err != nil {
			return err
		}
	}
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return err
	}
	if _, err := uuid.Parse(in.GetReviewerID()); err != nil {
		return err
	}
	if in.GetDomain() == "" {
		return fmt.Errorf("invalid domain")
	}
	if _, err := uuid.Parse(in.GetObjectID()); err != nil {
		return err
	}

	return nil
}

func ValidateManyCreate(in []*npool.ReviewReq) error {
	for _, info := range in {
		if err := ValidateCreate(info); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) CreateReview(ctx context.Context, in *npool.CreateReviewRequest) (*npool.CreateReviewResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateReview")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in.GetInfo())

	if err := ValidateCreate(in.GetInfo()); err != nil {
		return &npool.CreateReviewResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	span = commontracer.TraceInvoker(span, "review", "crud", "Create")

	info, err := crud.Create(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorf("fail create review: %v", err.Error())
		return &npool.CreateReviewResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateReviewResponse{
		Info: converter.Ent2Grpc(info),
	}, nil
}

func (s *Server) CreateReviews(ctx context.Context, in *npool.CreateReviewsRequest) (*npool.CreateReviewsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateReviews")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if len(in.GetInfos()) == 0 {
		return &npool.CreateReviewsResponse{}, status.Error(codes.InvalidArgument, "Infos is empty")
	}

	if err := ValidateManyCreate(in.GetInfos()); err != nil {
		return &npool.CreateReviewsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	span = tracer.TraceMany(span, in.GetInfos())
	span = commontracer.TraceInvoker(span, "review", "crud", "CreateBulk")

	rows, err := crud.CreateBulk(ctx, in.GetInfos())
	if err != nil {
		logger.Sugar().Errorf("fail create reviews: %v", err)
		return &npool.CreateReviewsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateReviewsResponse{
		Infos: converter.Ent2GrpcMany(rows),
	}, nil
}

func (s *Server) UpdateReview(ctx context.Context, in *npool.UpdateReviewRequest) (*npool.UpdateReviewResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateReview")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in.GetInfo())

	if _, err := uuid.Parse(in.GetInfo().GetID()); err != nil {
		logger.Sugar().Errorw("UpdateReview", "ID", in.GetInfo().GetID(), "error", err)
		return &npool.UpdateReviewResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if in.GetInfo().State != nil && in.GetInfo().GetState() == npool.ReviewState_Rejected {
		if in.GetInfo().GetMessage() == "" {
			logger.Sugar().Errorw("UpdateReview", "Message", in.GetInfo().GetMessage())
			return &npool.UpdateReviewResponse{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	span = commontracer.TraceInvoker(span, "review", "crud", "Update")

	info, err := crud.Update(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorf("fail create review: %v", err.Error())
		return &npool.UpdateReviewResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateReviewResponse{
		Info: converter.Ent2Grpc(info),
	}, nil
}

func (s *Server) GetReview(ctx context.Context, in *npool.GetReviewRequest) (*npool.GetReviewResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetReview")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceID(span, in.GetID())

	id, err := uuid.Parse(in.GetID())
	if err != nil {
		return &npool.GetReviewResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	span = commontracer.TraceInvoker(span, "review", "crud", "Row")

	info, err := crud.Row(ctx, id)
	if err != nil {
		logger.Sugar().Errorf("fail get review: %v", err)
		return &npool.GetReviewResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetReviewResponse{
		Info: converter.Ent2Grpc(info),
	}, nil
}

func ValidateConds(conds *npool.Conds) error { //nolint
	if conds.ID != nil {
		if _, err := uuid.Parse(conds.GetID().GetValue()); err != nil {
			return err
		}
	}
	if conds.AppID != nil {
		if _, err := uuid.Parse(conds.GetAppID().GetValue()); err != nil {
			return err
		}
	}
	if conds.ReviewerID != nil {
		if _, err := uuid.Parse(conds.GetReviewerID().GetValue()); err != nil {
			return err
		}
	}
	if conds.Domain != nil {
		if conds.GetDomain().GetValue() == "" {
			return fmt.Errorf("invalid domain")
		}
	}
	if conds.ObjectID != nil {
		if _, err := uuid.Parse(conds.GetObjectID().GetValue()); err != nil {
			return err
		}
	}
	if conds.Trigger != nil {
		switch npool.ReviewTriggerType(conds.GetTrigger().GetValue()) {
		case npool.ReviewTriggerType_AutoReviewed:
		case npool.ReviewTriggerType_LargeAmount:
		case npool.ReviewTriggerType_InsufficientFunds:
		case npool.ReviewTriggerType_InsufficientGas:
		case npool.ReviewTriggerType_InsufficientFundsGas:
		default:
			return fmt.Errorf("invalid trigger")
		}
	}
	if conds.ObjectType != nil {
		switch npool.ReviewObjectType(conds.GetObjectType().GetValue()) {
		case npool.ReviewObjectType_ObjectKyc:
		case npool.ReviewObjectType_ObjectWithdrawal:
		default:
			return fmt.Errorf("invalid object type")
		}
	}
	if conds.State != nil {
		switch npool.ReviewState(conds.GetState().GetValue()) {
		case npool.ReviewState_Wait:
		case npool.ReviewState_Approved:
		case npool.ReviewState_Rejected:
		default:
			return fmt.Errorf("invalid state")
		}
	}

	return nil
}

func (s *Server) GetReviewOnly(ctx context.Context, in *npool.GetReviewOnlyRequest) (*npool.GetReviewOnlyResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetReviewOnly")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.TraceConds(span, in.GetConds())

	if err := ValidateConds(in.GetConds()); err != nil {
		return &npool.GetReviewOnlyResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	span = commontracer.TraceInvoker(span, "review", "crud", "RowOnly")

	info, err := crud.RowOnly(ctx, in.GetConds())
	if err != nil {
		logger.Sugar().Errorf("fail get reviews: %v", err)
		return &npool.GetReviewOnlyResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetReviewOnlyResponse{
		Info: converter.Ent2Grpc(info),
	}, nil
}

func (s *Server) GetReviews(ctx context.Context, in *npool.GetReviewsRequest) (*npool.GetReviewsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetReviews")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.TraceConds(span, in.GetConds())
	span = commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))

	if err := ValidateConds(in.GetConds()); err != nil {
		return &npool.GetReviewsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	span = commontracer.TraceInvoker(span, "review", "crud", "Rows")

	rows, total, err := crud.Rows(ctx, in.GetConds(), int(in.GetOffset()), int(in.GetLimit()))
	if err != nil {
		logger.Sugar().Errorf("fail get reviews: %v", err)
		return &npool.GetReviewsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetReviewsResponse{
		Infos: converter.Ent2GrpcMany(rows),
		Total: uint32(total),
	}, nil
}

func (s *Server) ExistReview(ctx context.Context, in *npool.ExistReviewRequest) (*npool.ExistReviewResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "ExistReview")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceID(span, in.GetID())

	id, err := uuid.Parse(in.GetID())
	if err != nil {
		return &npool.ExistReviewResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	span = commontracer.TraceInvoker(span, "review", "crud", "Exist")

	exist, err := crud.Exist(ctx, id)
	if err != nil {
		logger.Sugar().Errorf("fail check review: %v", err)
		return &npool.ExistReviewResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.ExistReviewResponse{
		Info: exist,
	}, nil
}

func (s *Server) ExistReviewConds(ctx context.Context,
	in *npool.ExistReviewCondsRequest) (*npool.ExistReviewCondsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "ExistReviewConds")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.TraceConds(span, in.GetConds())

	if err := ValidateConds(in.GetConds()); err != nil {
		return &npool.ExistReviewCondsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	span = commontracer.TraceInvoker(span, "review", "crud", "ExistConds")

	exist, err := crud.ExistConds(ctx, in.GetConds())
	if err != nil {
		logger.Sugar().Errorf("fail check review: %v", err)
		return &npool.ExistReviewCondsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.ExistReviewCondsResponse{
		Info: exist,
	}, nil
}

func (s *Server) CountReviews(ctx context.Context, in *npool.CountReviewsRequest) (*npool.CountReviewsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CountReviews")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.TraceConds(span, in.GetConds())

	if err := ValidateConds(in.GetConds()); err != nil {
		return &npool.CountReviewsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	span = commontracer.TraceInvoker(span, "review", "crud", "Count")

	total, err := crud.Count(ctx, in.GetConds())
	if err != nil {
		logger.Sugar().Errorf("fail count reviews: %v", err)
		return &npool.CountReviewsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CountReviewsResponse{
		Info: total,
	}, nil
}

func (s *Server) DeleteReview(ctx context.Context, in *npool.DeleteReviewRequest) (*npool.DeleteReviewResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "DeleteReview")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceID(span, in.GetID())

	id, err := uuid.Parse(in.GetID())
	if err != nil {
		return &npool.DeleteReviewResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	span = commontracer.TraceInvoker(span, "review", "crud", "Delete")

	info, err := crud.Delete(ctx, id)
	if err != nil {
		logger.Sugar().Errorf("fail delete review: %v", err)
		return &npool.DeleteReviewResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteReviewResponse{
		Info: converter.Ent2Grpc(info),
	}, nil
}
