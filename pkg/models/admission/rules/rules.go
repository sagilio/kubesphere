package rules

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	admission "kubesphere.io/api/admission/v1alpha1"
	"kubesphere.io/kubesphere/pkg/api/admission/v1alpha1"
	kubesphere "kubesphere.io/kubesphere/pkg/client/clientset/versioned"
	"kubesphere.io/kubesphere/pkg/client/informers/externalversions"
	"kubesphere.io/kubesphere/pkg/models/admission/provider"
	"strings"
)

const (
	Kind = "rule.kubesphere.io"
)

type RuleManagerInterface interface {
	// GetRule gets the admission rule for the policy.
	GetRule(ctx context.Context, namespace, policyName, ruleName string) (*v1alpha1.RuleDetail, error)
	// ListRules lists the admission rules from the given policy.
	ListRules(ctx context.Context, namespace, policyName string) (*v1alpha1.RuleList, error)
	// CreateRule creates an admission rule for the policy.
	CreateRule(ctx context.Context, namespace, policyName string, rule *v1alpha1.PostRule) error
	// UpdateRule updates the admission rule for the policy with the given name.
	UpdateRule(ctx context.Context, namespace, policyName, ruleName string, rule *v1alpha1.PostRule) error
	// DeleteRule deletes the admission rule for the policy with the given name.
	DeleteRule(ctx context.Context, namespace, policyName, ruleName string) error
}

type RuleManager struct {
	ksClient    kubesphere.Interface
	ksInformers externalversions.SharedInformerFactory
	Providers   map[string]provider.Provider
}

func NewRuleManager(ksClient kubesphere.Interface, ksInformers externalversions.SharedInformerFactory, providers map[string]provider.Provider) *RuleManager {
	return &RuleManager{ksClient: ksClient, ksInformers: ksInformers, Providers: providers}
}

func (r RuleManager) GetRule(ctx context.Context, namespace, policyName, ruleName string) (*v1alpha1.RuleDetail, error) {
	panic("implement me")
}

func (r RuleManager) ListRules(ctx context.Context, namespace, policyName string) (*v1alpha1.RuleList, error) {
	panic("implement me")
}

func (r RuleManager) CreateRule(ctx context.Context, namespace, policyName string, rule *v1alpha1.PostRule) error {
	panic("implement me")
}

func (r RuleManager) UpdateRule(ctx context.Context, namespace, policyName, ruleName string, rule *v1alpha1.PostRule) error {
	panic("implement me")
}

func (r RuleManager) DeleteRule(ctx context.Context, namespace, policyName, ruleName string) error {
	panic("implement me")
}

func Rule(postRule *v1alpha1.PostRule, state admission.RuleState, namespaces ...string) *admission.Rule {
	if state == "" {
		state = admission.RuleInactive
	}
	rule := &admission.Rule{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(postRule.Name),
		},
		Spec: admission.RuleSpec{
			Name:        postRule.Name,
			Policy:      postRule.Policy,
			Provider:    postRule.Provider,
			Description: postRule.Description,
			Match: admission.Match{
				Namespaces: namespaces,
			},
			Parameters: runtime.RawExtension{
				Raw: []byte(postRule.Parameters),
			},
		},
		Status: admission.RuleStatus{State: state},
	}
	return rule
}
