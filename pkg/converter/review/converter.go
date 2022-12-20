package review

import (
	npool "github.com/NpoolPlatform/message/npool/review/mgr/v2"
	"github.com/NpoolPlatform/review-manager/pkg/db/ent"
)

func Ent2Grpc(row *ent.Review) *npool.Review {
	if row == nil {
		return nil
	}

	return &npool.Review{
		ID:         row.ID.String(),
		AppID:      row.AppID.String(),
		ReviewerID: row.ReviewerID.String(),
		Domain:     row.Domain,
		ObjectID:   row.ObjectID.String(),
		Trigger:    npool.ReviewTriggerType(npool.ReviewTriggerType_value[row.Trigger]),
		ObjectType: npool.ReviewObjectType(npool.ReviewObjectType_value[row.ObjectType]),
		State:      npool.ReviewState(npool.ReviewState_value[row.State]),
		Message:    row.Message,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}

func Ent2GrpcMany(rows []*ent.Review) []*npool.Review {
	infos := []*npool.Review{}
	for _, row := range rows {
		infos = append(infos, Ent2Grpc(row))
	}
	return infos
}
