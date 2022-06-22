package k8s

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	"github.com/yndd/ztp-dhcp/pkg/dhcp/testutils"
	"github.com/yndd/ztp-dhcp/pkg/structs"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ZtpK8sBackend struct {
	k8sclient runtimeclient.Client
}

// NewZtpK8sBackend returns a new Kubernetes based backend instance
func NewZtpK8sBackend(kubeconfig string) *ZtpK8sBackend {
	log.Infof("Instantiating K8sBackend")

	var config *rest.Config
	var err error

	// figure out where to load the kubeconfig from
	if kubeconfig == "" {
		log.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// create a scheme
	scheme := runtime.NewScheme()
	// add the required crd types to the scheme
	err = topov1alpha1.AddToScheme(scheme)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	// create a client including the scheme
	k8sclient, err := runtimeclient.New(config, runtimeclient.Options{
		Scheme: scheme,
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// finally return the backend instance
	return &ZtpK8sBackend{
		k8sclient: k8sclient,
	}
}

func (k *ZtpK8sBackend) GetDeviceInformation(cir *structs.ClientIdentifier) (*structs.DeviceInformation, error) {
	nl := &topov1alpha1.NodeList{}

	opts := runtimeclient.ListOptions{
		// TODO: Figure how to deal with namespaces
		Namespace: "default",
		Limit:     50,
	}

	var found = false
	var moreResults = true
	var result_node *topov1alpha1.Node = nil

	// get a pointer to the function that checks the specific CIType
	check_func := getCITypeCheckFunction(cir.CIType)

	// if not found and more results are available, continue to search
	for !found && moreResults {
		opts.Continue = nl.Continue
		// issue the kube-apiserver request
		err := k.k8sclient.List(context.TODO(), nl, &opts)
		if err != nil {
			return nil, fmt.Errorf("error fetching information from kubernets api. error: %v", err)
		}
		for _, entry := range nl.Items {
			// check if the actual entry is the one we are looking for
			if check_func(&entry, cir.Value) {
				result_node = &entry
				found = true
				break
			}
		}
		// we have more results if continue is set to something other then the emtpy string
		moreResults = nl.Continue != ""
	}

	// stop if the node could not be found
	if result_node == nil {
		return nil, fmt.Errorf("node with identifier %s, %s not found", cir.CIType.String(), cir.Value)
	}

	// populate the result object
	result := &structs.DeviceInformation{
		Name:              result_node.Name,
		MacAddress:        result_node.Spec.Properties.MacAddress,
		SerialNumber:      result_node.Spec.Properties.SerialNumber,
		CIDR:              result_node.Spec.Properties.MgmtIPAddress,
		Platform:          result_node.Spec.Properties.Platform,
		VendorType:        result_node.Spec.Properties.VendorType,
		Gateway:           "",
		ExpectedSWVersion: result_node.Spec.Properties.ExpectedSWVersion,
		NtpServersV4:      []string{},
		DnsServersV4:      []string{},
		Config:            "",
		Option66:          "",
		Option67:          "",
		Option43:          "",
	}
	return result, nil
}

// getCITypeCheckFunction returns a function that checks the given citype
func getCITypeCheckFunction(citype structs.CITypeEnum) func(*topov1alpha1.Node, string) bool {
	switch citype {
	case structs.MAC:
		return func(n *topov1alpha1.Node, mac string) bool {
			log.Debugf("mac check on %s ('%s' == '%s' => %s)", n.Name, n.Spec.Properties.MacAddress, mac, testutils.Bool2String(n.Spec.Properties.MacAddress == mac))
			return n.Spec.Properties.MacAddress == mac
		}
	case structs.String:
		return func(n *topov1alpha1.Node, serial string) bool {
			log.Debugf("serial check on %s ('%s' == '%s' => %s)", n.Name, n.Spec.Properties.SerialNumber, serial, testutils.Bool2String(n.Spec.Properties.MacAddress == serial))
			return n.Spec.Properties.SerialNumber == serial
		}
	}
	return nil
}
