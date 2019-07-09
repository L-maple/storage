package controller

import (
	"fmt"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type Controller struct {
	kubeClient      kubernetes.Interface
	podLister       corelisters.PodLister
	queue           workqueue.RateLimitingInterface
	podListerSynced cache.InformerSynced
}

func NewController(kubeCli kubernetes.Interface, informerFactory informers.SharedInformerFactory) *Controller {
	podInformer := informerFactory.Core().V1().Pods()
	c := &Controller{
		kubeClient: kubeCli,
		queue:      workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addPod,
		UpdateFunc: c.updatePod,
		DeleteFunc: c.deletePod,
	})
	c.podLister = podInformer.Lister()
	c.podListerSynced = podInformer.Informer().HasSynced
	return c
}

func (c *Controller) Run(workers int, stopCh <-chan struct{}) {
	defer c.queue.ShutDown()
	glog.Info("start pod controller")
	defer glog.Info("shutdown pod controller")
	//sync with pod
	if !cache.WaitForCacheSync(stopCh, c.podListerSynced) {
		return
	}
	for i := 0; i < workers; i++ {
		go wait.Until(c.worker, time.Second, stopCh)
	}
}

func (c *Controller) worker() {
	for c.processNextWrokItem() {
	}
}

func (c *Controller) processNextWrokItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	if err := c.sync(key.(string)); err != nil {
		c.queue.AddRateLimited(key)
	} else {
		c.queue.Forget(key)
	}
	return true
}

func (c *Controller) sync(key string) error {
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	pod, err := c.podLister.Pods(ns).Get(name)
	fmt.Printf("pod is %+v", *pod)
	return nil
}

func (c *Controller) checkResolve(pod *corev1.Pod) bool {
	annotations := pod.GetAnnotations()
	if status := annotations["scale.io/enable"]; status == "enable" {
		return true
	}
	return false
}

//只有需要存储的pod才会进入限速队列中
func (c *Controller) addPod(obj interface{}) {
	pod := obj.(*corev1.Pod)
	enable := c.checkResolve(pod)
	if enable {
		c.enqueuePod(pod)
	}
}

func (c *Controller) deletePod(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			glog.Errorf("can't get object form tombstone %+v", obj)
			return
		}
		pod, ok = tombstone.Obj.(*corev1.Pod)
		if !ok {
			glog.Errorf("tombstone contained object that is not a pod %+v", obj)
			return
		}
	}
	if enable := c.checkResolve(pod); enable {
		c.enqueuePod(pod)
	}
}

func (c *Controller) updatePod(old, cur interface{}) {
	curPod := cur.(*corev1.Pod)
	oldPod := old.(*corev1.Pod)
	if curPod.ResourceVersion == oldPod.ResourceVersion {
		return
	}
	if enable := c.checkResolve(curPod); enable {
		c.enqueuePod(curPod)
	}
}

func (c *Controller) enqueuePod(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		glog.Errorf("get key error, obj=%+v", obj)
		return
	}
	c.queue.Add(key)
}
