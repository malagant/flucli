package k8s

import (
	"context"
	"fmt"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Client wraps Kubernetes client functionality for FluxCD resources
type Client struct {
	client.Client
	kubernetes.Interface
	Config    *rest.Config
	Context   string
	Cluster   string
	Namespace string
}

// NewClient creates a new Kubernetes client
func NewClient(kubeconfig, context, namespace string) (*Client, error) {
	config, err := buildConfig(kubeconfig, context)
	if err != nil {
		return nil, fmt.Errorf("failed to build config: %w", err)
	}

	// Create controller-runtime client for CRDs
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add core/v1 to scheme: %w", err)
	}
	if err := sourcev1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add source/v1 to scheme: %w", err)
	}
	if err := sourcev1beta2.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add source/v1beta2 to scheme: %w", err)
	}
	if err := kustomizev1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add kustomize/v1 to scheme: %w", err)
	}
	if err := helmv2.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add helm/v2beta1 to scheme: %w", err)
	}

	ctrlClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("failed to create controller-runtime client: %w", err)
	}

	// Create standard Kubernetes client
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &Client{
		Client:    ctrlClient,
		Interface: k8sClient,
		Config:    config,
		Context:   context,
		Namespace: namespace,
	}, nil
}

// buildConfig builds a Kubernetes client configuration
func buildConfig(kubeconfig, context string) (*rest.Config, error) {
	configLoader := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeconfig != "" {
		configLoader.ExplicitPath = kubeconfig
	}

	configOverrides := &clientcmd.ConfigOverrides{}
	if context != "" {
		configOverrides.CurrentContext = context
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		configLoader,
		configOverrides,
	)

	return clientConfig.ClientConfig()
}

// TestConnection tests the connection to the Kubernetes cluster
func (c *Client) TestConnection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.CoreV1().Namespaces().Get(ctx, "default", metav1.GetOptions{})
	return err
}

// GetCurrentContext returns the current Kubernetes context
func (c *Client) GetCurrentContext() string {
	return c.Context
}

// GetClusterInfo returns basic cluster information
func (c *Client) GetClusterInfo(ctx context.Context) (*ClusterInfo, error) {
	version, err := c.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: %w", err)
	}

	nodes, err := c.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	return &ClusterInfo{
		Version:   version.String(),
		NodeCount: len(nodes.Items),
		Context:   c.Context,
	}, nil
}

// ClusterInfo contains basic cluster information
type ClusterInfo struct {
	Version   string
	NodeCount int
	Context   string
}
