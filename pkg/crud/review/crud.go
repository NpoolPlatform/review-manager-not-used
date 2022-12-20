package review

import (
	"context"
	"fmt"
	"time"

	constant "github.com/NpoolPlatform/review-manager/pkg/message/const"
	commontracer "github.com/NpoolPlatform/review-manager/pkg/tracer"
	tracer "github.com/NpoolPlatform/review-manager/pkg/tracer/review"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/review/mgr/v2"

	"github.com/NpoolPlatform/review-manager/pkg/db"
	"github.com/NpoolPlatform/review-manager/pkg/db/ent"
	"github.com/NpoolPlatform/review-manager/pkg/db/ent/review"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/google/uuid"
)

func CreateSet(c *ent.ReviewCreate, in *npool.ReviewReq) *ent.ReviewCreate {
	if in.ID != nil {
		c.SetID(uuid.MustParse(in.GetID()))
	}
	if in.AppID != nil {
		c.SetAppID(uuid.MustParse(in.GetAppID()))
	}
	if in.ReviewerID != nil {
		c.SetReviewerID(uuid.MustParse(in.GetReviewerID()))
	}
	if in.Domain != nil {
		c.SetDomain(in.GetDomain())
	}
	if in.ObjectID != nil {
		c.SetObjectID(uuid.MustParse(in.GetObjectID()))
	}
	if in.Trigger != nil {
		c.SetTrigger(in.GetTrigger().String())
	}
	if in.ObjectType != nil {
		c.SetObjectType(in.GetObjectType().String())
	}
	c.SetState(npool.ReviewState_Wait.String())
	return c
}

func Create(ctx context.Context, in *npool.ReviewReq) (*ent.Review, error) {
	var info *ent.Review
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "Create")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "db operation fail")
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in)

	err = db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		c := CreateSet(cli.Review.Create(), in)
		info, err = c.Save(_ctx)
		return err
	})
	if err != nil {
		return nil, err
	}

	return info, nil
}

func CreateBulk(ctx context.Context, in []*npool.ReviewReq) ([]*ent.Review, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateBulk")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "db operation fail")
			span.RecordError(err)
		}
	}()

	span = tracer.TraceMany(span, in)

	rows := []*ent.Review{}
	err = db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		bulk := make([]*ent.ReviewCreate, len(in))
		for i, info := range in {
			bulk[i] = CreateSet(tx.Review.Create(), info)
		}
		rows, err = tx.Review.CreateBulk(bulk...).Save(_ctx)
		return err
	})
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func UpdateSet(info *ent.Review, in *npool.ReviewReq) *ent.ReviewUpdateOne {
	stm := info.Update()

	if in.State != nil {
		stm = stm.SetState(in.GetState().String())
	}
	if in.Message != nil {
		stm = stm.SetMessage(in.GetMessage())
	}

	return stm
}

func Update(ctx context.Context, in *npool.ReviewReq) (*ent.Review, error) {
	var info *ent.Review
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "Update")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "db operation fail")
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in)

	err = db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		info, err = cli.Review.Query().Where(review.ID(uuid.MustParse(in.GetID()))).ForUpdate().Only(_ctx)
		if err != nil {
			return fmt.Errorf("fail query review: %v", err)
		}

		c := UpdateSet(info, in)
		info, err = c.Save(_ctx)
		return err
	})
	if err != nil {
		return nil, err
	}

	return info, nil
}

func Row(ctx context.Context, id uuid.UUID) (*ent.Review, error) {
	var info *ent.Review
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "Row")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "db operation fail")
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceID(span, id.String())

	err = db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		info, err = cli.Review.Query().Where(review.ID(id)).Only(_ctx)
		return err
	})
	if err != nil {
		return nil, err
	}

	return info, nil
}

func SetQueryConds(conds *npool.Conds, cli *ent.Client) (*ent.ReviewQuery, error) { //nolint
	stm := cli.Review.Query()
	if conds.ID != nil {
		switch conds.GetID().GetOp() {
		case cruder.EQ:
			stm.Where(review.ID(uuid.MustParse(conds.GetID().GetValue())))
		default:
			return nil, fmt.Errorf("invalid review field")
		}
	}
	if conds.AppID != nil {
		switch conds.GetAppID().GetOp() {
		case cruder.EQ:
			stm.Where(review.AppID(uuid.MustParse(conds.GetAppID().GetValue())))
		default:
			return nil, fmt.Errorf("invalid review field")
		}
	}
	if conds.ReviewerID != nil {
		switch conds.GetReviewerID().GetOp() {
		case cruder.EQ:
			stm.Where(review.ReviewerID(uuid.MustParse(conds.GetReviewerID().GetValue())))
		default:
			return nil, fmt.Errorf("invalid review field")
		}
	}
	if conds.Domain != nil {
		switch conds.GetDomain().GetOp() {
		case cruder.EQ:
			stm.Where(review.Domain(conds.GetDomain().GetValue()))
		default:
			return nil, fmt.Errorf("invalid review field")
		}
	}
	if conds.ObjectID != nil {
		switch conds.GetObjectID().GetOp() {
		case cruder.EQ:
			stm.Where(review.ObjectID(uuid.MustParse(conds.GetObjectID().GetValue())))
		default:
			return nil, fmt.Errorf("invalid review field")
		}
	}
	if conds.Trigger != nil {
		switch conds.GetTrigger().GetOp() {
		case cruder.EQ:
			stm.Where(review.Trigger(npool.ReviewTriggerType(conds.GetTrigger().GetValue()).String()))
		default:
			return nil, fmt.Errorf("invalid review field")
		}
	}
	if conds.ObjectType != nil {
		switch conds.GetObjectType().GetOp() {
		case cruder.EQ:
			stm.Where(review.ObjectType(npool.ReviewObjectType(conds.GetObjectType().GetValue()).String()))
		default:
			return nil, fmt.Errorf("invalid review field")
		}
	}
	if conds.State != nil {
		switch conds.GetState().GetOp() {
		case cruder.EQ:
			stm.Where(review.State(npool.ReviewState(conds.GetState().GetValue()).String()))
		default:
			return nil, fmt.Errorf("invalid review field")
		}
	}
	return stm, nil
}

func Rows(ctx context.Context, conds *npool.Conds, offset, limit int) ([]*ent.Review, int, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "Rows")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "db operation fail")
			span.RecordError(err)
		}
	}()

	span = tracer.TraceConds(span, conds)
	span = commontracer.TraceOffsetLimit(span, offset, limit)

	rows := []*ent.Review{}
	var total int
	err = db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		stm, err := SetQueryConds(conds, cli)
		if err != nil {
			return err
		}

		total, err = stm.Count(_ctx)
		if err != nil {
			return err
		}

		rows, err = stm.
			Offset(offset).
			Order(ent.Desc(review.FieldUpdatedAt)).
			Limit(limit).
			All(_ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func RowOnly(ctx context.Context, conds *npool.Conds) (*ent.Review, error) {
	var info *ent.Review
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "RowOnly")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "db operation fail")
			span.RecordError(err)
		}
	}()

	span = tracer.TraceConds(span, conds)

	err = db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		stm, err := SetQueryConds(conds, cli)
		if err != nil {
			return err
		}

		info, err = stm.Only(_ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return info, nil
}

func Count(ctx context.Context, conds *npool.Conds) (uint32, error) {
	var err error
	var total int

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "Count")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "db operation fail")
			span.RecordError(err)
		}
	}()

	span = tracer.TraceConds(span, conds)

	err = db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		stm, err := SetQueryConds(conds, cli)
		if err != nil {
			return err
		}

		total, err = stm.Count(_ctx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return uint32(total), nil
}

func Exist(ctx context.Context, id uuid.UUID) (bool, error) {
	var err error
	exist := false

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "Exist")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "db operation fail")
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceID(span, id.String())

	err = db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		exist, err = cli.Review.Query().Where(review.ID(id)).Exist(_ctx)
		return err
	})
	if err != nil {
		return false, err
	}

	return exist, nil
}

func ExistConds(ctx context.Context, conds *npool.Conds) (bool, error) {
	var err error
	exist := false

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "ExistConds")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "db operation fail")
			span.RecordError(err)
		}
	}()

	span = tracer.TraceConds(span, conds)

	err = db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		stm, err := SetQueryConds(conds, cli)
		if err != nil {
			return err
		}

		exist, err = stm.Exist(_ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	return exist, nil
}

func Delete(ctx context.Context, id uuid.UUID) (*ent.Review, error) {
	var info *ent.Review
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "Delete")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, "db operation fail")
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceID(span, id.String())

	err = db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		info, err = cli.Review.UpdateOneID(id).
			SetDeletedAt(uint32(time.Now().Unix())).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, err
	}

	return info, nil
}
