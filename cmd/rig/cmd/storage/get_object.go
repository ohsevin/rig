package storage

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rigdev/rig-go-api/api/v1/storage"
	"github.com/rigdev/rig-go-sdk"
	"github.com/rigdev/rig/cmd/rig/cmd/utils"
	"github.com/rigdev/rig/pkg/errors"
	"github.com/spf13/cobra"
)

func StorageGetObject(ctx context.Context, cmd *cobra.Command, args []string, nc rig.Client) error {
	var path string
	var err error
	if len(args) < 1 {
		path, err = utils.PromptGetInput("Object path:", utils.ValidateNonEmpty)
		if err != nil {
			return err
		}
	} else {
		path = args[0]
	}
	if isNSUri(path) {
		bucket, prefix, err := parseNSUri(path)
		if err != nil {
			return err
		}
		res, err := nc.Storage().GetObject(ctx, &connect.Request[storage.GetObjectRequest]{
			Msg: &storage.GetObjectRequest{
				Bucket: bucket,
				Path:   prefix,
			},
		})
		if err != nil {
			return err
		}

		if outputJson {
			cmd.Println(utils.ProtoToPrettyJson(res.Msg.GetObject()))
			return nil
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"Attribute", "Value"})
		t.AppendRows([]table.Row{
			{"Name", res.Msg.GetObject().GetPath()},
			{"Content type", res.Msg.GetObject().GetContentType()},
			{"Etag", res.Msg.GetObject().GetEtag()},
			{"Size", res.Msg.GetObject().GetSize()},
			{"Uploaded at", res.Msg.GetObject().GetLastModified().AsTime().Format("2006-01-02 15:04:05")},
		})
		cmd.Println(t.Render())

	} else {
		return errors.InvalidArgumentErrorf("invalid path: %s", path)
	}
	return nil
}