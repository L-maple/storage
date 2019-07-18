package rpc

import (
	"context"
	"errors"
	"github.com/golang/glog"
	"github.com/tommenx/cdproto/base"
	"github.com/tommenx/cdproto/cdpb"
	"google.golang.org/grpc"
)

var cli cdpb.CoordinatorClient

func Init(address string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		glog.Errorf("create grpc conn error, err=%+v", err)
		panic(err)
	}
	cli = cdpb.NewCoordinatorClient(conn)
}

func PutNodeStorage(ctx context.Context, node string, storage []*cdpb.Storage) error {
	req := &cdpb.PutNodeStorageRequest{
		Base: &base.Base{},
		Node: &cdpb.NodeStorage{},
	}
	req.Node.NodeName = node
	req.Node.Storage = storage
	rsp, err := cli.PutNodeStorage(ctx, req)
	if err != nil {
		glog.Errorf("call PutNodeStorage error, err=%+v", err)
		return err
	}
	if rsp.BaseResp.Code != 0 {
		glog.Errorf("remote server error,code=%v,msg=%v", rsp.BaseResp.Code, rsp.BaseResp.Message)
		return errors.New("remote server error")
	}
	return nil
}

func GetNodeStorage(ctx context.Context) (map[string]*cdpb.NodeStorage, error) {
	req := &cdpb.GetNodeStorageRequest{
		Base: &base.Base{},
	}
	rsp, err := cli.GetNodeStorage(ctx, req)
	if err != nil {
		glog.Errorf("call get node storage error, err=%+v", err)
		return nil, err
	}
	if rsp.BaseResp.Code != 0 {
		glog.Errorf("remote server error, coede=%v, msg=%+v", rsp.BaseResp.Code, rsp.BaseResp.Message)
		return nil, errors.New("remote server get node storage error")
	}
	return rsp.NodeMap, nil
}

func PutPodResource(ctx context.Context, ns, name string, request map[string]int64) error {
	req := &cdpb.PutPodResourceRequest{}
	pod := &cdpb.PodResource{}
	pod.Name = name
	pod.Namespace = ns
	pod.RequestResource = request
	req.Pod = pod
	rsp, err := cli.PutPodResource(ctx, req)
	if err != nil {
		glog.Errorf("call put pod resource error, err=%+v", err)
		return err
	}
	if rsp.BaseResp.Code != 0 {
		glog.Errorf("remote server error, code=%d, msg=%v", rsp.BaseResp.Code, rsp.BaseResp.Message)
		return errors.New("remote server put pod resource error")
	}
	return nil
}

func GetPodResource(ctx context.Context, ns, name string) (*cdpb.PodResource, error) {
	req := &cdpb.GetPodResourceRequest{
		Namespace: ns,
		Pod:       name,
	}
	rsp, err := cli.GetPodResource(ctx, req)
	if err != nil {
		glog.Errorf("call get pod resource error, err=%+v", err)
		return nil, err
	}
	if rsp.BaseResp.Code != 0 {
		glog.Errorf("remote server error, code=%d, msg=%v", rsp.BaseResp.Code, rsp.BaseResp.Message)
		return nil, errors.New("remote server get pod resource error")
	}
	return rsp.Pod, nil
}

func PutVolume(ctx context.Context, ns, pvc string, volume *cdpb.Volume) error {
	req := &cdpb.PutVolumeRequest{
		Base:      &base.Base{},
		Volume:    volume,
		Namespace: ns,
		Pvc:       pvc,
	}
	rsp, err := cli.PutVolume(ctx, req)
	if err != nil {
		glog.Errorf("call put volume error, err=%+v", err)
		return err
	}
	if rsp.BaseResp.Code != 0 {
		glog.Errorf("remote server error, code=%d, msg=%v", rsp.BaseResp.Code, rsp.BaseResp.Message)
		return errors.New("remote server put volume error")
	}
	return nil

}

func GetVolume(ctx context.Context, ns, pvc string) (*cdpb.Volume, error) {
	req := &cdpb.GetVolumeRequest{
		Base:      &base.Base{},
		Namespace: ns,
		Name:      pvc,
	}
	rsp, err := cli.GetVolume(ctx, req)
	if err != nil {
		glog.Errorf("call get volume error, err=%+v", err)
		return nil, err
	}
	if rsp.BaseResp.Code != 0 {
		glog.Errorf("remote server error, code=%d, msg=%v", rsp.BaseResp.Code, rsp.BaseResp.Message)
		return nil, errors.New("remote server get volume error")
	}
	return rsp.Volume, nil
}