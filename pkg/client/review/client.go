//nolint:dupl
package review

import (
	"context"
	"fmt"
	"time"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/review/mgr/v2"

	constant "github.com/NpoolPlatform/review-manager/pkg/message/const"
)

var timeout = 10 * time.Second

type handler func(context.Context, npool.ManagerClient) (cruder.Any, error)

func withCRUD(ctx context.Context, handler handler) (cruder.Any, error) {
	_ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	conn, err := grpc2.GetGRPCConn(constant.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get review connection: %v", err)
	}

	defer conn.Close()

	cli := npool.NewManagerClient(conn)

	return handler(_ctx, cli)
}

func CreateReview(ctx context.Context, in *npool.ReviewReq) (*npool.Review, error) {
	info, err := withCRUD(ctx, func(_ctx context.Context, cli npool.ManagerClient) (cruder.Any, error) {
		resp, err := cli.CreateReview(ctx, &npool.CreateReviewRequest{
			Info: in,
		})
		if err != nil {
			return nil, fmt.Errorf("fail create review: %v", err)
		}
		return resp.GetInfo(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail create review: %v", err)
	}
	return info.(*npool.Review), nil
}

func CreateReviews(ctx context.Context, in []*npool.ReviewReq) ([]*npool.Review, error) {
	infos, err := withCRUD(ctx, func(_ctx context.Context, cli npool.ManagerClient) (cruder.Any, error) {
		resp, err := cli.CreateReviews(ctx, &npool.CreateReviewsRequest{
			Infos: in,
		})
		if err != nil {
			return nil, fmt.Errorf("fail create reviews: %v", err)
		}
		return resp.GetInfos(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail create reviews: %v", err)
	}
	return infos.([]*npool.Review), nil
}

func UpdateReview(ctx context.Context, in *npool.ReviewReq) (*npool.Review, error) {
	info, err := withCRUD(ctx, func(_ctx context.Context, cli npool.ManagerClient) (cruder.Any, error) {
		resp, err := cli.UpdateReview(ctx, &npool.UpdateReviewRequest{
			Info: in,
		})
		if err != nil {
			return nil, fmt.Errorf("fail create review: %v", err)
		}
		return resp.GetInfo(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail create review: %v", err)
	}
	return info.(*npool.Review), nil
}
func GetReview(ctx context.Context, id string) (*npool.Review, error) {
	info, err := withCRUD(ctx, func(_ctx context.Context, cli npool.ManagerClient) (cruder.Any, error) {
		resp, err := cli.GetReview(ctx, &npool.GetReviewRequest{
			ID: id,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get review: %v", err)
		}
		return resp.GetInfo(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get review: %v", err)
	}
	return info.(*npool.Review), nil
}

func GetReviewOnly(ctx context.Context, conds *npool.Conds) (*npool.Review, error) {
	info, err := withCRUD(ctx, func(_ctx context.Context, cli npool.ManagerClient) (cruder.Any, error) {
		resp, err := cli.GetReviewOnly(ctx, &npool.GetReviewOnlyRequest{
			Conds: conds,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get review: %v", err)
		}
		return resp.GetInfo(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get review: %v", err)
	}
	return info.(*npool.Review), nil
}

func GetReviews(ctx context.Context, conds *npool.Conds, limit, offset int32) ([]*npool.Review, uint32, error) {
	var total uint32
	infos, err := withCRUD(ctx, func(_ctx context.Context, cli npool.ManagerClient) (cruder.Any, error) {
		resp, err := cli.GetReviews(ctx, &npool.GetReviewsRequest{
			Conds:  conds,
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get reviews: %v", err)
		}
		total = resp.GetTotal()
		return resp.GetInfos(), nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fail get reviews: %v", err)
	}
	return infos.([]*npool.Review), total, nil
}

func ExistReview(ctx context.Context, id string) (bool, error) {
	infos, err := withCRUD(ctx, func(_ctx context.Context, cli npool.ManagerClient) (cruder.Any, error) {
		resp, err := cli.ExistReview(ctx, &npool.ExistReviewRequest{
			ID: id,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get review: %v", err)
		}
		return resp.GetInfo(), nil
	})
	if err != nil {
		return false, fmt.Errorf("fail get review: %v", err)
	}
	return infos.(bool), nil
}

func ExistReviewConds(ctx context.Context, conds *npool.Conds) (bool, error) {
	infos, err := withCRUD(ctx, func(_ctx context.Context, cli npool.ManagerClient) (cruder.Any, error) {
		resp, err := cli.ExistReviewConds(ctx, &npool.ExistReviewCondsRequest{
			Conds: conds,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get review: %v", err)
		}
		return resp.GetInfo(), nil
	})
	if err != nil {
		return false, fmt.Errorf("fail get review: %v", err)
	}
	return infos.(bool), nil
}

func CountReviews(ctx context.Context, conds *npool.Conds) (uint32, error) {
	infos, err := withCRUD(ctx, func(_ctx context.Context, cli npool.ManagerClient) (cruder.Any, error) {
		resp, err := cli.CountReviews(ctx, &npool.CountReviewsRequest{
			Conds: conds,
		})
		if err != nil {
			return nil, fmt.Errorf("fail count review: %v", err)
		}
		return resp.GetInfo(), nil
	})
	if err != nil {
		return 0, fmt.Errorf("fail count review: %v", err)
	}
	return infos.(uint32), nil
}
