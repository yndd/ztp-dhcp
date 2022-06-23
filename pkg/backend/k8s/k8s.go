package k8s

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	"github.com/yndd/ztp-dhcp/pkg/dhcp/testutils"
	"github.com/yndd/ztp-dhcp/pkg/structs"
	"github.com/yndd/ztp-dhcp/pkg/utils"
	corev1 "k8s.io/api/core/v1"
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

	// add core types to scheme
	err = corev1.AddToScheme(scheme)
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

func (k *ZtpK8sBackend) GetWebserverInformation() (*structs.WebserverInfo, error) {
	serviceName := "ztp-webserver-yndd"
	protocolLookupLabel := "ztp.webserver.protocol"
	// figure out which namespace we are runnning in
	namespace := utils.DeduceNamespace("ndd-system")

	s := &corev1.Service{}
	err := k.k8sclient.Get(context.TODO(), runtimeclient.ObjectKey{
		Namespace: namespace,
		Name:      serviceName,
	}, s)

	if err != nil {
		return nil, fmt.Errorf("error retrieving ztp-webserver service information from k8s api: %v", err)
	}

	wsi := &structs.WebserverInfo{}

	switch s.Spec.Type {
	case corev1.ServiceTypeLoadBalancer:
		if s.Status.LoadBalancer.Ingress[0].Hostname != "" {
			wsi.IpFqdn = s.Status.LoadBalancer.Ingress[0].Hostname
		} else if s.Status.LoadBalancer.Ingress[0].IP != "" {
			wsi.IpFqdn = s.Status.LoadBalancer.Ingress[0].IP
		}
	case corev1.ServiceTypeExternalName:
		wsi.IpFqdn = s.Spec.ExternalName
	case corev1.ServiceTypeClusterIP:
		if len(s.Spec.ExternalIPs) > 0 {
			wsi.IpFqdn = s.Spec.ExternalIPs[0]
		}
	}
	// check if we found any endpoint information
	if wsi.IpFqdn == "" {
		return nil, fmt.Errorf("unable to determine external ClusterIP, ExternalName or ExternalIP for service '%s' in namespace '%s'", serviceName, namespace)
	}

	wsi.Port = s.Spec.Ports[0].Port

	if val, exists := s.ObjectMeta.Labels[protocolLookupLabel]; exists {
		wsi.Protocol = val
	} else {
		return nil, fmt.Errorf("unable to determine protocol of service '%s', missing label '%s'", serviceName, protocolLookupLabel)
	}

	return wsi, nil
}

func (k *ZtpK8sBackend) GetDeviceInformation(cir *structs.ClientIdentifier) (*structs.DeviceInformation, error) {
	// figure out the namespace we are running in atm.
	// if "POD_NAMESPACE" is not set fall back to "ndd-system" as namespace
	namespace := utils.DeduceNamespace("ndd-system")

	opts := runtimeclient.ListOptions{
		Namespace: namespace,
		Limit:     50,
	}

	nl := &topov1alpha1.NodeList{}

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
