// Code generated by clubbygen.
// GENERATED FILE DO NOT EDIT
// +build clubby_strict

package ota

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cesanta.com/common/go/mgrpc"
	"cesanta.com/common/go/mgrpc/frame"
	"cesanta.com/common/go/ourjson"
	"cesanta.com/common/go/ourtrace"
	"github.com/cesanta/errors"
	"golang.org/x/net/trace"

	"github.com/cesanta/ucl"
	"github.com/cesanta/validate-json/schema"
	"github.com/golang/glog"
)

var _ = bytes.MinRead
var _ = fmt.Errorf
var emptyMessage = ourjson.RawMessage{}
var _ = ourtrace.New
var _ = trace.New

const ServiceID = "http://cesanta.com/mg_rpc/serviceOTA"

type ListSectionsResult struct {
	Section  *string `json:"section,omitempty"`
	Version  *string `json:"version,omitempty"`
	Writable *bool   `json:"writable,omitempty"`
}

type UpdateArgs struct {
	Blob           *string `json:"blob,omitempty"`
	Blob_type      *string `json:"blob_type,omitempty"`
	Blob_url       *string `json:"blob_url,omitempty"`
	Commit_timeout *int64  `json:"commit_timeout,omitempty"`
	Section        *string `json:"section,omitempty"`
	Version        *string `json:"version,omitempty"`
}

type Service interface {
	Commit(ctx context.Context) error
	ListSections(ctx context.Context) ([]ListSectionsResult, error)
	Revert(ctx context.Context) error
	Update(ctx context.Context, args *UpdateArgs) error
}

type Instance interface {
	Call(context.Context, string, *frame.Command) (*frame.Response, error)
}

type _validators struct {
	// This comment prevents gofmt from aligning types in the struct.
	ListSectionsResult *schema.Validator
	// This comment prevents gofmt from aligning types in the struct.
	UpdateArgs *schema.Validator
}

var (
	validators     *_validators
	validatorsOnce sync.Once
)

func initValidators() {
	validators = &_validators{}

	loader := schema.NewLoader()

	service, err := ucl.Parse(bytes.NewBuffer(_ServiceDefinition))
	if err != nil {
		panic(err)
	}
	// Patch up shortcuts to be proper schemas.
	for _, v := range service.(*ucl.Object).Find("methods").(*ucl.Object).Value {
		if s, ok := v.(*ucl.Object).Find("result").(*ucl.String); ok {
			for kk := range v.(*ucl.Object).Value {
				if kk.Value == "result" {
					v.(*ucl.Object).Value[kk] = &ucl.Object{
						Value: map[ucl.Key]ucl.Value{
							ucl.Key{Value: "type"}: s,
						},
					}
				}
			}
		}
		if v.(*ucl.Object).Find("args") == nil {
			continue
		}
		args := v.(*ucl.Object).Find("args").(*ucl.Object)
		for kk, vv := range args.Value {
			if s, ok := vv.(*ucl.String); ok {
				args.Value[kk] = &ucl.Object{
					Value: map[ucl.Key]ucl.Value{
						ucl.Key{Value: "type"}: s,
					},
				}
			}
		}
	}
	var s *ucl.Object
	_ = s // avoid unused var error
	validators.ListSectionsResult, err = schema.NewValidator(service.(*ucl.Object).Find("methods").(*ucl.Object).Find("ListSections").(*ucl.Object).Find("result"), loader)
	if err != nil {
		panic(err)
	}
	s = &ucl.Object{
		Value: map[ucl.Key]ucl.Value{
			ucl.Key{Value: "properties"}: service.(*ucl.Object).Find("methods").(*ucl.Object).Find("Update").(*ucl.Object).Find("args"),
			ucl.Key{Value: "type"}:       &ucl.String{Value: "object"},
		},
	}
	if req, found := service.(*ucl.Object).Find("methods").(*ucl.Object).Find("Update").(*ucl.Object).Lookup("required_args"); found {
		s.Value[ucl.Key{Value: "required"}] = req
	}
	validators.UpdateArgs, err = schema.NewValidator(s, loader)
	if err != nil {
		panic(err)
	}
}

func NewClient(i Instance, addr string) Service {
	validatorsOnce.Do(initValidators)
	return &_Client{i: i, addr: addr}
}

type _Client struct {
	i    Instance
	addr string
}

func (c *_Client) Commit(ctx context.Context) (err error) {
	cmd := &frame.Command{
		Cmd: "OTA.Commit",
	}
	resp, err := c.i.Call(ctx, c.addr, cmd)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.Status != 0 {
		return errors.Trace(&mgrpc.ErrorResponse{Status: resp.Status, Msg: resp.StatusMsg})
	}
	return nil
}

func (c *_Client) ListSections(ctx context.Context) (res []ListSectionsResult, err error) {
	cmd := &frame.Command{
		Cmd: "OTA.ListSections",
	}
	resp, err := c.i.Call(ctx, c.addr, cmd)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.Status != 0 {
		return nil, errors.Trace(&mgrpc.ErrorResponse{Status: resp.Status, Msg: resp.StatusMsg})
	}

	bb, err := resp.Response.MarshalJSON()
	if err != nil {
		glog.Errorf("Failed to marshal result as JSON: %+v", err)
	} else {
		rv, err := ucl.Parse(bytes.NewReader(bb))
		if err == nil {
			if err := validators.ListSectionsResult.Validate(rv); err != nil {
				glog.Warningf("Got invalid result for ListSections: %+v", err)
				return nil, errors.Annotatef(err, "invalid response for ListSections")
			}
		}
	}
	var r []ListSectionsResult
	err = resp.Response.UnmarshalInto(&r)
	if err != nil {
		return nil, errors.Annotatef(err, "unmarshaling response")
	}
	return r, nil
}

func (c *_Client) Revert(ctx context.Context) (err error) {
	cmd := &frame.Command{
		Cmd: "OTA.Revert",
	}
	resp, err := c.i.Call(ctx, c.addr, cmd)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.Status != 0 {
		return errors.Trace(&mgrpc.ErrorResponse{Status: resp.Status, Msg: resp.StatusMsg})
	}
	return nil
}

func (c *_Client) Update(ctx context.Context, args *UpdateArgs) (err error) {
	cmd := &frame.Command{
		Cmd: "OTA.Update",
	}

	cmd.Args = ourjson.DelayMarshaling(args)
	b, err := cmd.Args.MarshalJSON()
	if err != nil {
		glog.Errorf("Failed to marshal args as JSON: %+v", err)
	} else {
		v, err := ucl.Parse(bytes.NewReader(b))
		if err != nil {
			glog.Errorf("Failed to parse just serialized JSON value %q: %+v", string(b), err)
		} else {
			if err := validators.UpdateArgs.Validate(v); err != nil {
				glog.Warningf("Sending invalid args for Update: %+v", err)
				return errors.Annotatef(err, "invalid args for Update")
			}
		}
	}
	resp, err := c.i.Call(ctx, c.addr, cmd)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.Status != 0 {
		return errors.Trace(&mgrpc.ErrorResponse{Status: resp.Status, Msg: resp.StatusMsg})
	}
	return nil
}

//func RegisterService(i *clubby.Instance, impl Service) error {
//validatorsOnce.Do(initValidators)
//s := &_Server{impl}
//i.RegisterCommandHandler("OTA.Commit", s.Commit)
//i.RegisterCommandHandler("OTA.ListSections", s.ListSections)
//i.RegisterCommandHandler("OTA.Revert", s.Revert)
//i.RegisterCommandHandler("OTA.Update", s.Update)
//i.RegisterService(ServiceID, _ServiceDefinition)
//return nil
//}

type _Server struct {
	impl Service
}

func (s *_Server) Commit(ctx context.Context, src string, cmd *frame.Command) (interface{}, error) {
	return nil, s.impl.Commit(ctx)
}

func (s *_Server) ListSections(ctx context.Context, src string, cmd *frame.Command) (interface{}, error) {
	r, err := s.impl.ListSections(ctx)
	if err != nil {
		return nil, errors.Trace(err)
	}
	bb, err := json.Marshal(r)
	if err == nil {
		v, err := ucl.Parse(bytes.NewBuffer(bb))
		if err != nil {
			glog.Errorf("Failed to parse just serialized JSON value %q: %+v", string(bb), err)
		} else {
			if err := validators.ListSectionsResult.Validate(v); err != nil {
				glog.Warningf("Returned invalid response for ListSections: %+v", err)
				return nil, errors.Annotatef(err, "server generated invalid responce for ListSections")
			}
		}
	}
	return r, nil
}

func (s *_Server) Revert(ctx context.Context, src string, cmd *frame.Command) (interface{}, error) {
	return nil, s.impl.Revert(ctx)
}

func (s *_Server) Update(ctx context.Context, src string, cmd *frame.Command) (interface{}, error) {
	b, err := cmd.Args.MarshalJSON()
	if err != nil {
		glog.Errorf("Failed to marshal args as JSON: %+v", err)
	} else {
		if v, err := ucl.Parse(bytes.NewReader(b)); err != nil {
			glog.Errorf("Failed to parse valid JSON value %q: %+v", string(b), err)
		} else {
			if err := validators.UpdateArgs.Validate(v); err != nil {
				glog.Warningf("Got invalid args for Update: %+v", err)
				return nil, errors.Annotatef(err, "invalid args for Update")
			}
		}
	}
	var args UpdateArgs
	if len(cmd.Args) > 0 {
		if err := cmd.Args.UnmarshalInto(&args); err != nil {
			return nil, errors.Annotatef(err, "unmarshaling args")
		}
	}
	return nil, s.impl.Update(ctx, &args)
}

var _ServiceDefinition = json.RawMessage([]byte(`{
  "doc": "OTA service provides a way to update device's software.",
  "methods": {
    "Commit": {
      "doc": "Commit a previously initiated update."
    },
    "ListSections": {
      "doc": "Returns a list of components of the device's software. Each section is updated individually.",
      "result": {
        "items": {
          "properties": {
            "section": {
              "type": "string"
            },
            "version": {
              "type": "string"
            },
            "writable": {
              "type": "boolean"
            }
          },
          "type": "object"
        },
        "type": "array"
      }
    },
    "Revert": {
      "doc": "Revert a previously initiated update."
    },
    "Update": {
      "args": {
        "blob": {
          "doc": "Image as a string, if appropriate.",
          "type": "string"
        },
        "blob_type": {
          "doc": "Type of the blob. Valid values: manifest, zip.",
          "type": "string"
        },
        "blob_url": {
          "doc": "URL pointing to the image if it's too big to fit in the ` + "`" + `blob` + "`" + `.",
          "type": "string"
        },
        "commit_timeout": {
          "doc": "Normally update is committed if firmware init succeeds, If timeout is set and non-zero, the update will require an explicit commit. If the specified time expires without a commit, update is rolled back.",
          "type": "integer"
        },
        "section": {
          "doc": "Name of the section to update.",
          "type": "string"
        },
        "version": {
          "doc": "Optional version of the new image.",
          "type": "string"
        }
      },
      "doc": "Instructs the device to update a given section."
    }
  },
  "name": "OTA",
  "namespace": "http://cesanta.com/mg_rpc/service",
  "visibility": "private"
}`))
