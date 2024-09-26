package main

import (
	clientset "builder/pkg/client/generated/clientset/versioned"
	informer "builder/pkg/client/generated/informers/externalversions"
	"builder/pkg/controller"
	"builder/pkg/signals"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil)
	ctx := signals.SetupSignalHandler()
	logger := klog.FromContext(ctx)

	cfg, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		logger.Error(err, "Error building kubeconfig")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	// new k8s client from cfg
	k8sClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		logger.Error(err, "Error building kubernetes clientset")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	client, err := clientset.NewForConfig(cfg)
	if err != nil {
		logger.Error(err, "Error building kubernetes clientset")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	factory := informer.NewSharedInformerFactory(client, 0)

	controller := controller.NewController(ctx, k8sClient, client,
		factory.Image().V1().Images(),
		factory.Builder().V1().Builders())

	factory.Start(ctx.Done())

	if err = controller.Run(ctx, 2); err != nil {
		logger.Error(err, "Error running controller")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}
}
