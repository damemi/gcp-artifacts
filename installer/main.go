package main

import (
	"context"
	"fmt"
	"os"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/odigos-io/odigos/api/generated/odigos/clientset/versioned/typed/odigos/v1alpha1"
	"github.com/odigos-io/odigos/cli/cmd/resources"
	"github.com/odigos-io/odigos/cli/cmd/resources/resourcemanager"
	"github.com/odigos-io/odigos/cli/pkg/kube"
	"github.com/odigos-io/odigos/common"
	"github.com/odigos-io/odigos/k8sutils/pkg/installationmethod"
)

func main() {
	ctx := context.Background()

	fmt.Println("Starting Odigos installer")

	fmt.Println("Getting k8s config")
	k8sConfig, err := config.GetConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to get k8s config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Creating k8s clientset")
	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to create clientset: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Creating k8s dynamic client")
	dynamicClient, err := dynamic.NewForConfig(k8sConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to get k8s dynamic client: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Creating k8s extend clientset")
	extendClientset, err := apiextensionsclient.NewForConfig(k8sConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to get k8s extendClientset: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Creating Odigos client")
	odigosClient, err := v1alpha1.NewForConfig(k8sConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to get Odigos client: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Creating kube client")
	kubeClient := &kube.Client{
		Interface:     clientset,
		Clientset:     clientset,
		Dynamic:       dynamicClient,
		ApiExtensions: extendClientset,
		OdigosClient:  odigosClient,
		Config:        k8sConfig,
	}

	ns := os.Getenv("ODIGOS_NAMESPACE")
	onPremToken := os.Getenv("ODIGOS_ON_PREM_TOKEN")
	var odigosProToken string
	odigosTier := common.CommunityOdigosTier
	if onPremToken != "" {
		odigosTier = common.OnPremOdigosTier
		odigosProToken = onPremToken
	}
	odigosConfiguration := common.OdigosConfiguration{}
	version := os.Getenv("ODIGOS_VERSION")
	managerOpts := resourcemanager.ManagerOpts{}

	imageReferences := resourcemanager.ImageReferences{
		AutoscalerImage:    os.Getenv("ODIGOS_AUTOSCALER_IMAGE"),
		CollectorImage:     os.Getenv("ODIGOS_COLLECTOR_IMAGE"),
		InitContainerImage: os.Getenv("ODIGOS_INIT_CONTAINER_IMAGE"),
		SchedulerImage:     os.Getenv("ODIGOS_SCHEDULER_IMAGE"),
		InstrumentorImage:  os.Getenv("ODIGOS_INSTRUMENTOR_IMAGE"),
		OdigletImage:       os.Getenv("ODIGOS_ODIGLET_IMAGE"),
		UIImage:            os.Getenv("ODIGOS_UI_IMAGE"),
	}
	managerOpts.ImageReferences = imageReferences

	odigosInstallerName := os.Getenv("ODIGOS_INSTALLER_NAME")
	odigosInstallerNamespace := os.Getenv("ODIGOS_INSTALLER_NAMESPACE")

	fmt.Println("Getting installer deployment")
	if odigosInstallerName != "" && odigosInstallerNamespace != "" {
		deployment, err := clientset.AppsV1().Deployments(odigosInstallerNamespace).Get(ctx, odigosInstallerName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to get installer deployment %s in namespace %s: %v\n", odigosInstallerName, odigosInstallerNamespace, err)
			os.Exit(1)
		}

		isController := true
		blockOwnerDeletion := true
		ownerRef := metav1.OwnerReference{
			APIVersion:         "apps/v1",
			Kind:               "Deployment",
			Name:               deployment.Name,
			UID:                deployment.UID,
			Controller:         &isController,
			BlockOwnerDeletion: &blockOwnerDeletion,
		}
		managerOpts.OwnerReferences = []metav1.OwnerReference{ownerRef}
	}

	fmt.Println("Creating resource managers")
	resourceManagers := resources.CreateResourceManagers(
		kubeClient,
		ns,
		odigosTier,
		&odigosProToken,
		&odigosConfiguration,
		version,
		installationmethod.K8sInstallationMethodOdigosOperator,
		managerOpts)

	fmt.Println("Applying resource managers")
	err = resources.ApplyResourceManagers(ctx, kubeClient, resourceManagers, "Creating")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to apply resource managers: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Odigos installation completed successfully")
}
