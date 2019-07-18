package controller

import (
	"errors"
	"github.com/golang/glog"
	"k8s.io/client-go/informers"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

var (
	ErrNotFound = errors.New("pv not found")
	ErrNotBound = errors.New("pv not bounded")
)

type pvController struct {
	pvLister       corelisters.PersistentVolumeLister
	pvListerSynced cache.InformerSynced
}

type PVController interface {
	Run(stop <-chan struct{})
	GetPVCByPV(pvName string) (string, string, error)
}

func NewPVController(informerFactory informers.SharedInformerFactory) PVController {
	pvInformfer := informerFactory.Core().V1().PersistentVolumes()
	controller := &pvController{}
	controller.pvLister = pvInformfer.Lister()
	controller.pvListerSynced = pvInformfer.Informer().HasSynced
	return controller
}

func (c *pvController) Run(stopCh <-chan struct{}) {
	if !cache.WaitForCacheSync(stopCh, c.pvListerSynced) {
		glog.Error("sync pv timeout")
		return
	}
}

func (c *pvController) GetPVCByPV(pvName string) (string, string, error) {
	pv, err := c.pvLister.Get(pvName)
	if err != nil {
		glog.Errorf("get pv info error, pv name=%s, err=%+v", pvName, err)
		return "", "", err
	}
	if pv == nil {
		glog.Errorf("pv not found, pv name=%s", pvName)
		return "", "", ErrNotFound
	}
	if pv.Spec.ClaimRef == nil {
		glog.Errorf("can not found pv claim, pv name=%s", pvName)
		return "", "", ErrNotBound
	}
	return pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name, nil
}