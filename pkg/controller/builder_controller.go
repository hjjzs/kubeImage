package controller

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	"golang.org/x/time/rate"

	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	builderv1 "builder/pkg/apis/builder/v1"
	clientset "builder/pkg/client/generated/clientset/versioned"
	samplescheme "builder/pkg/client/generated/clientset/versioned/scheme"
	builderInformers "builder/pkg/client/generated/informers/externalversions/builder/v1"
	imageInformers "builder/pkg/client/generated/informers/externalversions/image/v1"

	buildListers "builder/pkg/client/generated/listers/builder/v1"
	imageListers "builder/pkg/client/generated/listers/image/v1"
)

const controllerAgentName = "builder-controller"

const (
	//// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	//SuccessSynced = "Synced"
	//// ErrResourceExists is used as part of the Event 'reason' when a Foo fails
	//// to sync due to a Deployment of the same name already existing.
	//ErrResourceExists = "ErrResourceExists"
	//
	//// MessageResourceExists is the message used for Events when a resource
	//// fails to sync due to a Deployment already existing
	//MessageResourceExists = "Resource %q already exists and is not managed by Foo"
	//// MessageResourceSynced is the message used for an Event fired when a Foo
	//// is synced successfully
	//MessageResourceSynced = "Foo synced successfully"
	//// FieldManager distinguishes this controller from other things writing to API objects
	//FieldManager = controllerAgentName

	ContextGetting      = "Getting"
	ImageBuilding       = "Building"
	ImagePushing        = "Pushing"
	ImageSourceCreating = "Creating"
	Finished            = "Finished"
	Failed              = "Failed"
)

// Controller is the controller implementation for Foo resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// sampleclientset is a clientset for our own API group
	client clientset.Interface

	imageList     imageListers.ImageLister
	imageSynced   cache.InformerSynced
	builderLister buildListers.BuilderLister
	builderSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.TypedRateLimitingInterface[cache.ObjectName]
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new sample controller
func NewController(
	ctx context.Context,
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	ImageInformer imageInformers.ImageInformer,
	BuilderInformer builderInformers.BuilderInformer) *Controller {
	logger := klog.FromContext(ctx)

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	utilruntime.Must(samplescheme.AddToScheme(scheme.Scheme))
	logger.V(4).Info("Creating event broadcaster")

	eventBroadcaster := record.NewBroadcaster(record.WithContext(ctx))
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})
	ratelimiter := workqueue.NewTypedMaxOfRateLimiter(
		workqueue.NewTypedItemExponentialFailureRateLimiter[cache.ObjectName](5*time.Millisecond, 1000*time.Second),
		&workqueue.TypedBucketRateLimiter[cache.ObjectName]{Limiter: rate.NewLimiter(rate.Limit(50), 300)},
	)

	controller := &Controller{
		kubeclientset: kubeclientset,
		client:        sampleclientset,
		builderLister: BuilderInformer.Lister(),
		builderSynced: BuilderInformer.Informer().HasSynced,
		imageList:     ImageInformer.Lister(),
		imageSynced:   ImageInformer.Informer().HasSynced,
		workqueue:     workqueue.NewTypedRateLimitingQueue(ratelimiter),
		recorder:      recorder,
	}

	logger.Info("Setting up event handlers")
	// Set up an event handler for when Foo resources change
	BuilderInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueFoo,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueFoo(new)
		},
		DeleteFunc: controller.enqueueFoo,
	})

	//deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
	//	AddFunc: controller.handleObject,
	//	UpdateFunc: func(old, new interface{}) {
	//		newDepl := new.(*appsv1.Deployment)
	//		oldDepl := old.(*appsv1.Deployment)
	//		if newDepl.ResourceVersion == oldDepl.ResourceVersion {
	//			// Periodic resync will send update events for all known Deployments.
	//			// Two different versions of the same Deployment will always have different RVs.
	//			return
	//		}
	//		controller.handleObject(new)
	//	},
	//	DeleteFunc: controller.handleObject,
	//})

	return controller
}

func (c *Controller) Run(ctx context.Context, workers int) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()
	logger := klog.FromContext(ctx)

	// Start the informer factories to begin populating the informer caches
	logger.Info("Starting Foo controller")

	// Wait for the caches to be synced before starting workers
	logger.Info("Waiting for informer caches to sync")

	if ok := cache.WaitForCacheSync(ctx.Done(), c.builderSynced, c.imageSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	logger.Info("Starting workers", "count", workers)
	// Launch two workers to process Foo resources
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.runWorker, time.Second)
	}

	logger.Info("Started workers")
	<-ctx.Done()
	logger.Info("Shutting down workers")

	return nil
}

func (c *Controller) enqueueFoo(obj interface{}) {
	if objectRef, err := cache.ObjectToName(obj); err != nil {
		utilruntime.HandleError(err)
		return
	} else {
		c.workqueue.Add(objectRef)
	}
}

func (c *Controller) runWorker(ctx context.Context) {
	for c.processNextWorkItem(ctx) {
	}
}

func (c *Controller) processNextWorkItem(ctx context.Context) bool {
	objRef, shutdown := c.workqueue.Get()
	logger := klog.FromContext(ctx)

	if shutdown {
		return false
	}

	// We call Done at the end of this func so the workqueue knows we have
	// finished processing this item. We also must remember to call Forget
	// if we do not want this work item being re-queued. For example, we do
	// not call Forget if a transient error occurs, instead the item is
	// put back on the workqueue and attempted again after a back-off
	// period.
	defer c.workqueue.Done(objRef)

	// Run the syncHandler, passing it the structured reference to the object to be synced.
	err := c.syncHandler(ctx, objRef)
	if err == nil {
		// If no error occurs then we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(objRef)
		logger.Info("Successfully synced", "objectName", objRef)
		return true
	}

	if c.workqueue.NumRequeues(objRef) < 3 {
		utilruntime.HandleErrorWithContext(ctx, err, "Error syncing; requeuing for later retry", "objectReference", objRef)
		c.workqueue.AddRateLimited(objRef)
		return true
	}

	utilruntime.HandleErrorWithContext(ctx, err, "Error syncing; dropping", "objectReference", objRef)
	get, _ := c.builderLister.Get(objRef.Name)
	get.Status.State = Failed
	c.client.BuilderV1().Builders().UpdateStatus(ctx, get, metav1.UpdateOptions{})
	c.workqueue.Forget(objRef)
	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Foo resource
// with the current status of the resource.
func (c *Controller) syncHandler(ctx context.Context, obj cache.ObjectName) error {
	logger := klog.LoggerWithValues(klog.FromContext(ctx), "objectRef", obj)

	builder, err := c.client.BuilderV1().Builders().Get(ctx, obj.Name, metav1.GetOptions{})
	if err != nil {
		logger.Info("start delete builder", "builder", obj.Name)
		return c.handlerDeleteBuilder(ctx, obj.Name)
	}

	logger.Info("start sync builder", "builder", obj.Name)

	switch builder.Status.State {
	case ContextGetting:
		err = c.handlerContextGetting(ctx, builder, logger)
	//case ImageBuilding:
	//	err = handerImageBuilding(ctx, c.client, builder)
	//case ImagePushing:
	//	err = handerImagePushing(ctx, c.client, builder)
	//case ImageSourceCreating:
	//	err = handerImageSourceCreating(ctx, c.client, builder)
	//case Finished:
	//	err = handerFinished(ctx, c.client, builder)
	case Failed:
		return nil
	default:
		err = c.handlerContextGetting(ctx, builder, logger)
	}

	return err
}

func (c *Controller) handlerContextGetting(ctx context.Context, builder *builderv1.Builder, logger klog.Logger) error {
	err := c.updateBuilderStatus(ctx, builder, ContextGetting)
	if err != nil {
		logger.Error(err, "update builder status failed")
		return err
	}
	// get downloader

	return nil
}

func (c *Controller) handlerDeleteBuilder(ctx context.Context, name string) error {

	return nil
}

func (c *Controller) updateBuilderStatus(ctx context.Context, builder *builderv1.Builder, status string) error {
	deepCopy := builder.DeepCopy()
	deepCopy.Status.State = status
	_, err := c.client.BuilderV1().Builders().UpdateStatus(ctx, builder, metav1.UpdateOptions{})
	return err
}

//func (c *Controller) handleObject(obj interface{}) {
//	var object metav1.Object
//	var ok bool
//	logger := klog.FromContext(context.Background())
//	if object, ok = obj.(metav1.Object); !ok {
//		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
//		if !ok {
//			// If the object value is not too big and does not contain sensitive information then
//			// it may be useful to include it.
//			utilruntime.HandleErrorWithContext(context.Background(), nil, "Error decoding object, invalid type", "type", fmt.Sprintf("%T", obj))
//			return
//		}
//		object, ok = tombstone.Obj.(metav1.Object)
//		if !ok {
//			// If the object value is not too big and does not contain sensitive information then
//			// it may be useful to include it.
//			utilruntime.HandleErrorWithContext(context.Background(), nil, "Error decoding object tombstone, invalid type", "type", fmt.Sprintf("%T", tombstone.Obj))
//			return
//		}
//		logger.V(4).Info("Recovered deleted object", "resourceName", object.GetName())
//	}
//	logger.V(4).Info("Processing object", "object", klog.KObj(object))
//	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
//		// If this object is not owned by a Foo, we should not do anything more
//		// with it.
//		if ownerRef.Kind != "Foo" {
//			return
//		}
//
//		foo, err := c.foosLister.Foos(object.GetNamespace()).Get(ownerRef.Name)
//		if err != nil {
//			logger.V(4).Info("Ignore orphaned object", "object", klog.KObj(object), "foo", ownerRef.Name)
//			return
//		}
//
//		c.enqueueFoo(foo)
//		return
//	}
//}
//
