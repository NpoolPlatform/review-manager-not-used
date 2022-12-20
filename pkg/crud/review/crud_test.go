package review

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/NpoolPlatform/review-manager/pkg/db/ent"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	valuedef "github.com/NpoolPlatform/message/npool"
	npool "github.com/NpoolPlatform/message/npool/review/mgr/v2"
	testinit "github.com/NpoolPlatform/review-manager/pkg/testinit"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func init() {
	if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); err == nil && runByGithubAction {
		return
	}
	if err := testinit.Init(); err != nil {
		fmt.Printf("cannot init test stub: %v\n", err)
	}
}

var ret = ent.Review{
	ID:         uuid.New(),
	AppID:      uuid.New(),
	ReviewerID: uuid.New(),
	Domain:     uuid.NewString(),
	ObjectID:   uuid.New(),
	Trigger:    npool.ReviewTriggerType_AutoReviewed.String(),
	ObjectType: npool.ReviewObjectType_ObjectKyc.String(),
	State:      npool.ReviewState_Wait.String(),
}

var (
	id         = ret.ID.String()
	appID      = ret.AppID.String()
	reviewerID = ret.ReviewerID.String()
	objectID   = ret.ObjectID.String()
	trigger    = npool.ReviewTriggerType(npool.ReviewTriggerType_value[ret.Trigger])
	objectType = npool.ReviewObjectType(npool.ReviewObjectType_value[ret.ObjectType])

	req = npool.ReviewReq{
		ID:         &id,
		AppID:      &appID,
		ReviewerID: &reviewerID,
		Domain:     &ret.Domain,
		ObjectID:   &objectID,
		Trigger:    &trigger,
		ObjectType: &objectType,
	}
)

var info *ent.Review

func create(t *testing.T) {
	var err error
	info, err = Create(context.Background(), &req)
	if assert.Nil(t, err) {
		ret.UpdatedAt = info.UpdatedAt
		ret.CreatedAt = info.CreatedAt
		assert.Equal(t, info.String(), ret.String())
	}
}

func createBulk(t *testing.T) {
	entities := []*ent.Review{
		{
			ID:         uuid.New(),
			AppID:      uuid.New(),
			ReviewerID: uuid.New(),
			Domain:     uuid.NewString(),
			ObjectID:   uuid.New(),
			Trigger:    npool.ReviewTriggerType_AutoReviewed.String(),
			ObjectType: npool.ReviewObjectType_ObjectKyc.String(),
			State:      npool.ReviewState_Wait.String(),
		},
		{
			ID:         uuid.New(),
			AppID:      uuid.New(),
			ReviewerID: uuid.New(),
			Domain:     uuid.NewString(),
			ObjectID:   uuid.New(),
			Trigger:    npool.ReviewTriggerType_LargeAmount.String(),
			ObjectType: npool.ReviewObjectType_ObjectWithdrawal.String(),
			State:      npool.ReviewState_Wait.String(),
		},
	}

	reqs := []*npool.ReviewReq{}
	for _, _ret := range entities {
		_id := _ret.ID.String()
		_appID := _ret.AppID.String()
		_reviewerID := _ret.ReviewerID.String()
		_objectID := _ret.ObjectID.String()
		_trigger := npool.ReviewTriggerType(npool.ReviewTriggerType_value[_ret.Trigger])
		_objectType := npool.ReviewObjectType(npool.ReviewObjectType_value[_ret.ObjectType])

		reqs = append(reqs, &npool.ReviewReq{
			ID:         &_id,
			AppID:      &_appID,
			ReviewerID: &_reviewerID,
			Domain:     &_ret.Domain,
			ObjectID:   &_objectID,
			Trigger:    &_trigger,
			ObjectType: &_objectType,
		})
	}
	infos, err := CreateBulk(context.Background(), reqs)
	if assert.Nil(t, err) {
		assert.Equal(t, len(infos), 2)
	}
}

func row(t *testing.T) {
	var err error
	info, err = Row(context.Background(), ret.ID)
	if assert.Nil(t, err) {
		assert.Equal(t, info.String(), ret.String())
	}
}

func update(t *testing.T) {
	var err error

	state := npool.ReviewState_Approved
	ret.State = state.String()
	req.State = &state

	info, err = Update(context.Background(), &req)
	if assert.Nil(t, err) {
		ret.UpdatedAt = info.UpdatedAt
		assert.Equal(t, info.String(), ret.String())
	}
}

func rows(t *testing.T) {
	infos, total, err := Rows(context.Background(),
		&npool.Conds{
			ID: &valuedef.StringVal{
				Value: id,
				Op:    cruder.EQ,
			},
		}, 0, 0)
	if assert.Nil(t, err) {
		if assert.Equal(t, total, 1) {
			assert.Equal(t, infos[0].String(), ret.String())
		}
	}
}

func rowOnly(t *testing.T) {
	var err error
	info, err = RowOnly(context.Background(),
		&npool.Conds{
			ID: &valuedef.StringVal{
				Value: id,
				Op:    cruder.EQ,
			},
		})
	if assert.Nil(t, err) {
		assert.Equal(t, info.String(), ret.String())
	}
}

func count(t *testing.T) {
	count, err := Count(context.Background(),
		&npool.Conds{
			ID: &valuedef.StringVal{
				Value: id,
				Op:    cruder.EQ,
			},
		},
	)
	if assert.Nil(t, err) {
		assert.Equal(t, count, uint32(1))
	}
}

func exist(t *testing.T) {
	exist, err := Exist(context.Background(), ret.ID)
	if assert.Nil(t, err) {
		assert.Equal(t, exist, true)
	}
}

func existConds(t *testing.T) {
	exist, err := ExistConds(context.Background(),
		&npool.Conds{
			ID: &valuedef.StringVal{
				Value: id,
				Op:    cruder.EQ,
			},
		},
	)
	if assert.Nil(t, err) {
		assert.Equal(t, exist, true)
	}
}

func deleteA(t *testing.T) {
	info, err := Delete(context.Background(), ret.ID)
	if assert.Nil(t, err) {
		ret.DeletedAt = info.DeletedAt
		assert.Equal(t, info.String(), ret.String())
	}
}

func TestReview(t *testing.T) {
	if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); err == nil && runByGithubAction {
		return
	}
	t.Run("create", create)
	t.Run("createBulk", createBulk)
	t.Run("update", update)
	t.Run("row", row)
	t.Run("rows", rows)
	t.Run("rowOnly", rowOnly)
	t.Run("exist", exist)
	t.Run("existConds", existConds)
	t.Run("count", count)
	t.Run("delete", deleteA)
}
