package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alpine-hodler/sherpa/internal/storage"
	"github.com/alpine-hodler/sherpa/pkg/proto"
	"github.com/alpine-hodler/sherpa/tools"
	"google.golang.org/protobuf/types/known/structpb"
)

type Generic interface {
	storage.Storage

	UpsertRawJSON(context.Context, *Raw, *proto.CreateResponse) error
}

type generic struct{ storage.Storage }

func New(ctx context.Context, dns string) (Generic, error) {
	stg, err := storage.New(ctx, dns)
	return &generic{stg}, err
}

// UpserRawJSON will upsert a Raw struct into the repository.
func (svc *generic) UpsertRawJSON(ctx context.Context, raw *Raw, rsp *proto.CreateResponse) error {
	var records []*structpb.Struct
	var data interface{}
	if err := json.Unmarshal(raw.Data, &data); err != nil {
		return fmt.Errorf("failed to unmarshal raw data: %w", err)
	}

	if err := tools.MakeRecordsRequest(data, &records); err != nil {
		return fmt.Errorf("error making records request: %v", err)
	}
	req := new(proto.UpsertRequest)
	req.Table = raw.Table
	req.Records = records
	return svc.Storage.Upsert(ctx, req, rsp)
}
