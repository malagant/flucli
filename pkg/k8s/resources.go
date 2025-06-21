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
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ResourceType represents the type of FluxCD resource
type ResourceType string

const (
	ResourceTypeGitRepository  ResourceType = "GitRepository"
	ResourceTypeHelmRepository ResourceType = "HelmRepository"
	ResourceTypeKustomization  ResourceType = "Kustomization"
	ResourceTypeHelmRelease    ResourceType = "HelmRelease"
)

// Resource represents a generic FluxCD resource
type Resource struct {
	Type        ResourceType  `json:"type"`
	Name        string        `json:"name"`
	Namespace   string        `json:"namespace"`
	Ready       bool          `json:"ready"`
	Status      string        `json:"status"`
	Message     string        `json:"message"`
	Age         time.Duration `json:"age"`
	LastUpdate  time.Time     `json:"last_update"`
	Conditions  []Condition   `json:"conditions"`
	Suspended   bool          `json:"suspended"`
	Source      string        `json:"source,omitempty"`
	Path        string        `json:"path,omitempty"`
	Revision    string        `json:"revision,omitempty"`
	URL         string        `json:"url,omitempty"`
	Chart       string        `json:"chart,omitempty"`
	Version     string        `json:"version,omitempty"`
}

// Condition represents a status condition
type Condition struct {
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	Reason             string    `json:"reason"`
	Message            string    `json:"message"`
	LastTransitionTime time.Time `json:"lastTransitionTime"`
}

// ListGitRepositories lists all GitRepository resources
func (c *Client) ListGitRepositories(ctx context.Context, namespace string) ([]Resource, error) {
	var gitRepos sourcev1.GitRepositoryList
	opts := []client.ListOption{}
	if namespace != "" {
		opts = append(opts, client.InNamespace(namespace))
	}

	if err := c.List(ctx, &gitRepos, opts...); err != nil {
		return nil, fmt.Errorf("failed to list GitRepositories: %w", err)
	}

	resources := make([]Resource, 0, len(gitRepos.Items))
	for _, repo := range gitRepos.Items {
		resource := Resource{
			Type:       ResourceTypeGitRepository,
			Name:       repo.Name,
			Namespace:  repo.Namespace,
			Age:        time.Since(repo.CreationTimestamp.Time),
			LastUpdate: time.Now(),
			Suspended:  repo.Spec.Suspend,
			URL:        repo.Spec.URL,
		}

		// Parse status
		if repo.Status.Conditions != nil {
			for _, cond := range repo.Status.Conditions {
				resource.Conditions = append(resource.Conditions, Condition{
					Type:               cond.Type,
					Status:             string(cond.Status),
					Reason:             cond.Reason,
					Message:            cond.Message,
					LastTransitionTime: cond.LastTransitionTime.Time,
				})

				if cond.Type == "Ready" {
					resource.Ready = cond.Status == metav1.ConditionTrue
					resource.Status = cond.Reason
					resource.Message = cond.Message
				}
			}
		}

		if repo.Status.Artifact != nil {
			resource.Revision = repo.Status.Artifact.Revision
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

// ListHelmRepositories lists all HelmRepository resources
func (c *Client) ListHelmRepositories(ctx context.Context, namespace string) ([]Resource, error) {
	var helmRepos sourcev1beta2.HelmRepositoryList
	opts := []client.ListOption{}
	if namespace != "" {
		opts = append(opts, client.InNamespace(namespace))
	}

	if err := c.List(ctx, &helmRepos, opts...); err != nil {
		return nil, fmt.Errorf("failed to list HelmRepositories: %w", err)
	}

	resources := make([]Resource, 0, len(helmRepos.Items))
	for _, repo := range helmRepos.Items {
		resource := Resource{
			Type:       ResourceTypeHelmRepository,
			Name:       repo.Name,
			Namespace:  repo.Namespace,
			Age:        time.Since(repo.CreationTimestamp.Time),
			LastUpdate: time.Now(),
			Suspended:  repo.Spec.Suspend,
			URL:        repo.Spec.URL,
		}

		// Parse status
		if repo.Status.Conditions != nil {
			for _, cond := range repo.Status.Conditions {
				resource.Conditions = append(resource.Conditions, Condition{
					Type:               cond.Type,
					Status:             string(cond.Status),
					Reason:             cond.Reason,
					Message:            cond.Message,
					LastTransitionTime: cond.LastTransitionTime.Time,
				})

				if cond.Type == "Ready" {
					resource.Ready = cond.Status == metav1.ConditionTrue
					resource.Status = cond.Reason
					resource.Message = cond.Message
				}
			}
		}

		if repo.Status.Artifact != nil {
			resource.Revision = repo.Status.Artifact.Revision
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

// ListKustomizations lists all Kustomization resources
func (c *Client) ListKustomizations(ctx context.Context, namespace string) ([]Resource, error) {
	var kustomizations kustomizev1.KustomizationList
	opts := []client.ListOption{}
	if namespace != "" {
		opts = append(opts, client.InNamespace(namespace))
	}

	if err := c.List(ctx, &kustomizations, opts...); err != nil {
		return nil, fmt.Errorf("failed to list Kustomizations: %w", err)
	}

	resources := make([]Resource, 0, len(kustomizations.Items))
	for _, ks := range kustomizations.Items {
		resource := Resource{
			Type:       ResourceTypeKustomization,
			Name:       ks.Name,
			Namespace:  ks.Namespace,
			Age:        time.Since(ks.CreationTimestamp.Time),
			LastUpdate: time.Now(),
			Suspended:  ks.Spec.Suspend,
			Path:       ks.Spec.Path,
		}

		if ks.Spec.SourceRef.Kind == "GitRepository" {
			resource.Source = ks.Spec.SourceRef.Name
		}

		// Parse status
		if ks.Status.Conditions != nil {
			for _, cond := range ks.Status.Conditions {
				resource.Conditions = append(resource.Conditions, Condition{
					Type:               cond.Type,
					Status:             string(cond.Status),
					Reason:             cond.Reason,
					Message:            cond.Message,
					LastTransitionTime: cond.LastTransitionTime.Time,
				})

				if cond.Type == "Ready" {
					resource.Ready = cond.Status == metav1.ConditionTrue
					resource.Status = cond.Reason
					resource.Message = cond.Message
				}
			}
		}

		if ks.Status.LastAppliedRevision != "" {
			resource.Revision = ks.Status.LastAppliedRevision
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

// ListHelmReleases lists all HelmRelease resources
func (c *Client) ListHelmReleases(ctx context.Context, namespace string) ([]Resource, error) {
	var helmReleases helmv2.HelmReleaseList
	opts := []client.ListOption{}
	if namespace != "" {
		opts = append(opts, client.InNamespace(namespace))
	}

	if err := c.List(ctx, &helmReleases, opts...); err != nil {
		return nil, fmt.Errorf("failed to list HelmReleases: %w", err)
	}

	resources := make([]Resource, 0, len(helmReleases.Items))
	for _, hr := range helmReleases.Items {
		resource := Resource{
			Type:       ResourceTypeHelmRelease,
			Name:       hr.Name,
			Namespace:  hr.Namespace,
			Age:        time.Since(hr.CreationTimestamp.Time),
			LastUpdate: time.Now(),
			Suspended:  hr.Spec.Suspend,
			Chart:      hr.Spec.Chart.Spec.Chart,
			Version:    hr.Spec.Chart.Spec.Version,
		}

		if hr.Spec.Chart.Spec.SourceRef.Kind == "HelmRepository" {
			resource.Source = hr.Spec.Chart.Spec.SourceRef.Name
		}

		// Parse status
		if hr.Status.Conditions != nil {
			for _, cond := range hr.Status.Conditions {
				resource.Conditions = append(resource.Conditions, Condition{
					Type:               cond.Type,
					Status:             string(cond.Status),
					Reason:             cond.Reason,
					Message:            cond.Message,
					LastTransitionTime: cond.LastTransitionTime.Time,
				})

				if cond.Type == "Ready" {
					resource.Ready = cond.Status == metav1.ConditionTrue
					resource.Status = cond.Reason
					resource.Message = cond.Message
				}
			}
		}

		if hr.Status.LastAppliedRevision != "" {
			resource.Revision = hr.Status.LastAppliedRevision
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

// SuspendResource suspends a FluxCD resource
func (c *Client) SuspendResource(ctx context.Context, resourceType ResourceType, name, namespace string) error {
	return c.updateSuspendStatus(ctx, resourceType, name, namespace, true)
}

// ResumeResource resumes a FluxCD resource
func (c *Client) ResumeResource(ctx context.Context, resourceType ResourceType, name, namespace string) error {
	return c.updateSuspendStatus(ctx, resourceType, name, namespace, false)
}

// updateSuspendStatus updates the suspend status of a resource
func (c *Client) updateSuspendStatus(ctx context.Context, resourceType ResourceType, name, namespace string, suspend bool) error {
	var obj client.Object
	
	switch resourceType {
	case ResourceTypeGitRepository:
		obj = &sourcev1.GitRepository{}
	case ResourceTypeHelmRepository:
		obj = &sourcev1beta2.HelmRepository{}
	case ResourceTypeKustomization:
		obj = &kustomizev1.Kustomization{}
	case ResourceTypeHelmRelease:
		obj = &helmv2.HelmRelease{}
	default:
		return fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	key := types.NamespacedName{Name: name, Namespace: namespace}
	if err := c.Get(ctx, key, obj); err != nil {
		return fmt.Errorf("failed to get %s/%s: %w", resourceType, name, err)
	}

	// Update suspend field based on resource type
	switch resourceType {
	case ResourceTypeGitRepository:
		repo := obj.(*sourcev1.GitRepository)
		repo.Spec.Suspend = suspend
	case ResourceTypeHelmRepository:
		repo := obj.(*sourcev1beta2.HelmRepository)
		repo.Spec.Suspend = suspend
	case ResourceTypeKustomization:
		ks := obj.(*kustomizev1.Kustomization)
		ks.Spec.Suspend = suspend
	case ResourceTypeHelmRelease:
		hr := obj.(*helmv2.HelmRelease)
		hr.Spec.Suspend = suspend
	}

	if err := c.Update(ctx, obj); err != nil {
		return fmt.Errorf("failed to update %s/%s: %w", resourceType, name, err)
	}

	return nil
}

// ReconcileResource triggers reconciliation of a FluxCD resource
func (c *Client) ReconcileResource(ctx context.Context, resourceType ResourceType, name, namespace string) error {
	var obj client.Object
	
	switch resourceType {
	case ResourceTypeGitRepository:
		obj = &sourcev1.GitRepository{}
	case ResourceTypeHelmRepository:
		obj = &sourcev1beta2.HelmRepository{}
	case ResourceTypeKustomization:
		obj = &kustomizev1.Kustomization{}
	case ResourceTypeHelmRelease:
		obj = &helmv2.HelmRelease{}
	default:
		return fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	key := types.NamespacedName{Name: name, Namespace: namespace}
	if err := c.Get(ctx, key, obj); err != nil {
		return fmt.Errorf("failed to get %s/%s: %w", resourceType, name, err)
	}

	// Add reconcile annotation
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations["reconcile.fluxcd.io/requestedAt"] = time.Now().UTC().Format(time.RFC3339)
	obj.SetAnnotations(annotations)

	if err := c.Update(ctx, obj); err != nil {
		return fmt.Errorf("failed to update %s/%s: %w", resourceType, name, err)
	}

	return nil
}

// GetEvents returns Kubernetes events related to FluxCD resources
func (c *Client) GetEvents(ctx context.Context, namespace string) ([]corev1.Event, error) {
	eventList, err := c.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{
		FieldSelector: "involvedObject.apiVersion=source.toolkit.fluxcd.io/v1,kustomize.toolkit.fluxcd.io/v1,helm.toolkit.fluxcd.io/v2beta1",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	return eventList.Items, nil
}
