package policies

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	admission "kubesphere.io/api/admission/v1alpha1"
	"kubesphere.io/kubesphere/pkg/api/admission/v1alpha1"
	kubesphere "kubesphere.io/kubesphere/pkg/client/clientset/versioned"
	"kubesphere.io/kubesphere/pkg/client/informers/externalversions"
	"kubesphere.io/kubesphere/pkg/models/admission/provider"
	"strings"
)

type PolicyManagerInterface interface {
	// GetPolicy gets the admission policy with the given name.
	GetPolicy(ctx context.Context, policyName string) (*v1alpha1.PolicyDetail, error)
	// ListPolicies lists the admission policies with the given name.
	ListPolicies(ctx context.Context) (*v1alpha1.PolicyList, error)
	// CreatePolicy creates an admission policy.
	CreatePolicy(ctx context.Context, policy *v1alpha1.PostPolicy) error
	// UpdatePolicy updates the admission policy with the given name.
	UpdatePolicy(ctx context.Context, policyName string, policy *v1alpha1.PostPolicy) error
	// DeletePolicy deletes the admission policy with the given name.
	DeletePolicy(ctx context.Context, policyName string) error
}

type PolicyManager struct {
	ksClient    kubesphere.Interface
	ksInformers externalversions.SharedInformerFactory
	Providers   map[string]provider.Provider
}

func NewPolicyManager(ksClient kubesphere.Interface, ksInformers externalversions.SharedInformerFactory, providers map[string]provider.Provider) *PolicyManager {
	return &PolicyManager{ksClient: ksClient, ksInformers: ksInformers, Providers: providers}
}

func (p PolicyManager) GetPolicy(_ context.Context, policyName string) (*v1alpha1.PolicyDetail, error) {
	policy, err := p.ksInformers.Admission().V1alpha1().Policies().Lister().Get(policyName)
	if err != nil {
		return nil, err
	}
	return PolicyDetail(policy), nil
}

func (p PolicyManager) ListPolicies(ctx context.Context) (*v1alpha1.PolicyList, error) {
	panic("implement me")
}

func (p PolicyManager) CreatePolicy(ctx context.Context, policy *v1alpha1.PostPolicy) error {
	panic("implement me")
}

func (p PolicyManager) UpdatePolicy(ctx context.Context, policyName string, policy *v1alpha1.PostPolicy) error {
	panic("implement me")
}

func (p PolicyManager) DeletePolicy(ctx context.Context, policyName string) error {
	panic("implement me")
}

func PolicyDetail(policy *admission.Policy) *v1alpha1.PolicyDetail {
	panic("implement me")
}

func Policy(postPolicy *v1alpha1.PostPolicy, state admission.PolicyState) *admission.Policy {
	if state == "" {
		state = admission.PolicyInactive
	}

	params := postPolicy.Parameters
	targets := postPolicy.Targets

	var policyTargets []admission.PolicyContentTarget
	for _, target := range targets {
		policyTargets = append(policyTargets, admission.PolicyContentTarget{
			Target:     target.Target,
			Expression: target.Expression,
			Import:     target.Import,
		})
	}

	policy := &admission.Policy{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(postPolicy.Name),
		},
		Spec: admission.PolicySpec{
			Name:           postPolicy.Name,
			PolicyTemplate: postPolicy.PolicyTemplate,
			Description:    postPolicy.Description,
			Provider:       postPolicy.Provider,
			Content: admission.PolicyContent{
				Spec: admission.PolicyContentSpec{
					Names: admission.Names{
						Name: postPolicy.Name,
					},
					Parameters: admission.Parameters{
						Validation: &admission.Validation{
							OpenAPIV3Schema: params.Validation.OpenAPIV3Schema,
							LegacySchema:    params.Validation.LegacySchema,
						},
					},
				},
				Targets: policyTargets,
			},
		},
		Status: admission.PolicyStatus{State: state},
	}
	return policy
}

func PolicyFromTemplate(template *admission.PolicyTemplate, name string, provider string, state admission.PolicyState) *admission.Policy {
	if state == "" {
		state = admission.PolicyInactive
	}
	if name == "" {
		name = template.Name
	}
	templateContent := template.Spec.Content
	targets := templateContent.Targets
	var policyTargets []admission.PolicyContentTarget
	for _, target := range targets {
		if target.Provider != provider {
			continue
		}
		policyTargets = append(policyTargets, admission.PolicyContentTarget{
			Target:     target.Target,
			Expression: target.Expression,
			Import:     target.Import,
		})
	}

	policy := &admission.Policy{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(name),
		},
		Spec: admission.PolicySpec{
			Name:           name,
			PolicyTemplate: template.Spec.Name,
			Description:    template.Spec.Description,
			Provider:       provider,
			Content: admission.PolicyContent{
				Spec: admission.PolicyContentSpec{
					Names:      templateContent.Spec.Names,
					Parameters: templateContent.Spec.Parameters,
				},
				Targets: policyTargets,
			},
		},
		Status: admission.PolicyStatus{State: state},
	}
	return policy
}
