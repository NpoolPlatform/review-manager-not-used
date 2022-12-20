package review

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	trace1 "go.opentelemetry.io/otel/trace"

	npool "github.com/NpoolPlatform/message/npool/review/mgr/v2"
)

func trace(span trace1.Span, in *npool.ReviewReq, index int) trace1.Span {
	span.SetAttributes(
		attribute.String(fmt.Sprintf("ID.%v", index), in.GetID()),
		attribute.String(fmt.Sprintf("AppID.%v", index), in.GetAppID()),
		attribute.String(fmt.Sprintf("ReviewerID.%v", index), in.GetReviewerID()),
		attribute.String(fmt.Sprintf("Domain.%v", index), in.GetDomain()),
		attribute.String(fmt.Sprintf("ObjectID.%v", index), in.GetObjectID()),
		attribute.String(fmt.Sprintf("Trigger.%v", index), in.GetTrigger().String()),
		attribute.String(fmt.Sprintf("ObjectType.%v", index), in.GetObjectType().String()),
		attribute.String(fmt.Sprintf("State.%v", index), in.GetState().String()),
	)
	return span
}

func Trace(span trace1.Span, in *npool.ReviewReq) trace1.Span {
	return trace(span, in, 0)
}

func TraceConds(span trace1.Span, in *npool.Conds) trace1.Span {
	span.SetAttributes(
		attribute.String("ID.Op", in.GetID().GetOp()),
		attribute.String("ID.Value", in.GetID().GetValue()),
		attribute.String("AppID.Op", in.GetAppID().GetOp()),
		attribute.String("AppID.Value", in.GetAppID().GetValue()),
		attribute.String("ReviewerID.Op", in.GetReviewerID().GetOp()),
		attribute.String("ReviewerID.Value", in.GetReviewerID().GetValue()),
		attribute.String("Domain.Op", in.GetDomain().GetOp()),
		attribute.String("Domain.Value", in.GetDomain().GetValue()),
		attribute.String("ObjectID.Op", in.GetObjectID().GetOp()),
		attribute.String("ObjectID.Value", in.GetObjectID().GetValue()),
		attribute.String("Trigger.Op", in.GetTrigger().GetOp()),
		attribute.String("Trigger.Value", npool.ReviewTriggerType(in.GetTrigger().GetValue()).String()),
		attribute.String("ObjectType.Op", in.GetObjectType().GetOp()),
		attribute.String("ObjectType.Value", npool.ReviewObjectType(in.GetObjectType().GetValue()).String()),
		attribute.String("State.Op", in.GetState().GetOp()),
		attribute.String("State.Value", npool.ReviewState(in.GetState().GetValue()).String()),
	)
	return span
}

func TraceMany(span trace1.Span, infos []*npool.ReviewReq) trace1.Span {
	for index, info := range infos {
		span = trace(span, info, index)
	}
	return span
}
