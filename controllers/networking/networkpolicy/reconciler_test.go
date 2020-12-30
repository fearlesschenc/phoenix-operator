package networkpolicy

import (
	"context"
	"github.com/fearlesschenc/kubesphere/pkg/constants"
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
)

const SampleNetworkPolicyName = "networkpolicy-sample"

func IsNetworkPolicyAllowFrom(obj *networkingv1.NetworkPolicy, namespace string) bool {
	for _, peer := range obj.Spec.Ingress[0].From {
		selector, err := metav1.LabelSelectorAsSelector(peer.NamespaceSelector)
		Expect(err).NotTo(HaveOccurred())

		if selector.Matches(labels.Set{constants.NamespaceLabelKey: namespace}) {
			return true
		}
	}

	return false
}

var _ = Describe("Network Policy Reconcile Sanity", func() {
	ctx := context.TODO()

	SetupManager(ctx)

	var policy *networkingv1alpha1.NetworkPolicy
	JustBeforeEach(func() {
		err := k8sClient.Create(ctx, policy)
		Expect(err).NotTo(HaveOccurred())
	})
	// clean up
	JustAfterEach(func() {
		err := k8sClient.Delete(ctx, policy)
		Expect(err).NotTo(HaveOccurred())

		policies := &networkingv1.NetworkPolicyList{}
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, types.NamespacedName{Name: policy.Name}, policy); !errors.IsNotFound(err) {
				return false
			}

			if err := k8sClient.List(ctx, policies); err != nil {
				return false
			}

			return len(policies.Items) == 0
		}, "5s").Should(BeTrue())
	})

	When("specify network policy accept from workspace to workspace", func() {
		BeforeEach(func() {
			policy = &networkingv1alpha1.NetworkPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name: SampleNetworkPolicyName,
				},
				Spec: networkingv1alpha1.NetworkPolicySpec{
					Workspace:         "foo",
					NamespaceSelector: metav1.LabelSelector{},
					From: []networkingv1alpha1.NetworkPolicyPeer{
						{
							Workspace:         "bar",
							NamespaceSelector: metav1.LabelSelector{},
						},
					},
				},
			}
		})

		It("should generate corresponding network policy", func() {
			Eventually(func() bool {
				obj := &networkingv1.NetworkPolicy{}

				for _, namespace := range workspaces[0].Namespaces {
					objKey := types.NamespacedName{Namespace: namespace.Name, Name: SampleNetworkPolicyName}
					if err := k8sClient.Get(ctx, objKey, obj); err != nil {
						return false
					}
				}

				return true
			}, "5s").Should(BeTrue())
		})

		It("should have corresponding namespace in network policy", func() {
			for _, namespace := range workspaces[0].Namespaces {
				obj := &networkingv1.NetworkPolicy{}
				objKey := types.NamespacedName{Namespace: namespace.Name, Name: SampleNetworkPolicyName}
				err := k8sClient.Get(ctx, objKey, obj)
				Expect(err).NotTo(HaveOccurred())

				Expect(IsNetworkPolicyAllowFrom(obj, "bar1")).To(BeTrue())
				Expect(IsNetworkPolicyAllowFrom(obj, "bar2")).To(BeTrue())
			}
		})
	})

	//Specify("workspace-application", func() {
	//	policy := &networkingv1alpha1.NetworkPolicy{
	//		ObjectMeta: metav1.ObjectMeta{
	//			Name: SampleNetworkPolicyName,
	//		},
	//		Spec: networkingv1alpha1.NetworkPolicySpec{
	//			Workspace:         "foo",
	//			NamespaceSelector: metav1.LabelSelector{},
	//			From: []networkingv1alpha1.NetworkPolicyPeer{
	//				{
	//					Workspace:         "bar",
	//					NamespaceSelector: metav1.LabelSelector{MatchLabels: map[string]string{constants.NamespaceLabelKey: "bar1"}},
	//				},
	//			},
	//		},
	//	}
	//
	//	err := k8sClient.Create(context.TODO(), policy)
	//	Expect(err).NotTo(HaveOccurred())
	//
	//	Specify("namespace should create correspond networkpolicy", func() {
	//		Eventually(func() bool {
	//			obj := &networkingv1.NetworkPolicy{}
	//
	//			for _, namespace := range workspaces[0].Namespaces {
	//				objKey := types.NamespacedName{Namespace: namespace.Name, Name: SampleNetworkPolicyName}
	//				if err := k8sClient.Get(context.TODO(), objKey, obj); err != nil {
	//					return false
	//				}
	//			}
	//
	//			return true
	//		}).Should(BeTrue())
	//	})
	//
	//	Specify("network policy should allow correspond namespace access", func() {
	//		for _, namespace := range workspaces[0].Namespaces {
	//			obj := &networkingv1.NetworkPolicy{}
	//			objKey := types.NamespacedName{Namespace: namespace.Name, Name: SampleNetworkPolicyName}
	//			err := k8sClient.Get(context.TODO(), objKey, obj)
	//			Expect(err).NotTo(HaveOccurred())
	//
	//			Expect(IsNetworkPolicyAllowFrom(obj, "bar1")).To(BeTrue())
	//			Expect(IsNetworkPolicyAllowFrom(obj, "bar2")).To(BeFalse())
	//		}
	//	})
	//})
	//
	//Specify("application-workspace", func() {
	//	policy := &networkingv1alpha1.NetworkPolicy{
	//		ObjectMeta: metav1.ObjectMeta{
	//			Name: SampleNetworkPolicyName,
	//		},
	//		Spec: networkingv1alpha1.NetworkPolicySpec{
	//			Workspace:         "foo",
	//			NamespaceSelector: metav1.LabelSelector{MatchLabels: map[string]string{constants.NamespaceLabelKey: "foo1"}},
	//			From: []networkingv1alpha1.NetworkPolicyPeer{
	//				{
	//					Workspace: "bar",
	//				},
	//			},
	//		},
	//	}
	//
	//	err := k8sClient.Create(context.TODO(), policy)
	//	Expect(err).NotTo(HaveOccurred())
	//
	//	Specify("namespace should create correspond networkpolicy", func() {
	//		Eventually(func() bool {
	//			obj := &networkingv1.NetworkPolicy{}
	//
	//			objKey := types.NamespacedName{Namespace: "foo1", Name: SampleNetworkPolicyName}
	//			if err := k8sClient.Get(context.TODO(), objKey, obj); err != nil {
	//				return false
	//			}
	//
	//			return true
	//		}).Should(BeTrue())
	//
	//		Consistently(func() bool {
	//			obj := &networkingv1.NetworkPolicy{}
	//
	//			objKey := types.NamespacedName{Namespace: "foo2", Name: SampleNetworkPolicyName}
	//			if err := k8sClient.Get(context.TODO(), objKey, obj); err != nil {
	//				return false
	//			}
	//
	//			return true
	//
	//		}).Should(BeFalse())
	//	})
	//
	//	Specify("network policy should allow correspond namespace access", func() {
	//		obj := &networkingv1.NetworkPolicy{}
	//		objKey := types.NamespacedName{Namespace: "foo1", Name: SampleNetworkPolicyName}
	//		err := k8sClient.Get(context.TODO(), objKey, obj)
	//		Expect(err).NotTo(HaveOccurred())
	//		for _, namespace := range workspaces[1].Namespaces {
	//			Expect(IsNetworkPolicyAllowFrom(obj, namespace.Name)).To(BeTrue())
	//		}
	//	})
	//})
	//
	//Specify("application-application", func() {
	//	policy := &networkingv1alpha1.NetworkPolicy{
	//		ObjectMeta: metav1.ObjectMeta{
	//			Name: SampleNetworkPolicyName,
	//		},
	//		Spec: networkingv1alpha1.NetworkPolicySpec{
	//			Workspace:         "foo",
	//			NamespaceSelector: metav1.LabelSelector{MatchLabels: map[string]string{constants.NamespaceLabelKey: "foo1"}},
	//			From: []networkingv1alpha1.NetworkPolicyPeer{
	//				{
	//					Workspace:         "bar",
	//					NamespaceSelector: metav1.LabelSelector{MatchLabels: map[string]string{constants.NamespaceLabelKey: "bar1"}},
	//				},
	//			},
	//		},
	//	}
	//
	//	err := k8sClient.Create(context.TODO(), policy)
	//	Expect(err).NotTo(HaveOccurred())
	//
	//	Specify("namespace should create correspond networkpolicy", func() {
	//		Eventually(func() bool {
	//			obj := &networkingv1.NetworkPolicy{}
	//
	//			objKey := types.NamespacedName{Namespace: "foo1", Name: SampleNetworkPolicyName}
	//			if err := k8sClient.Get(context.TODO(), objKey, obj); err != nil {
	//				return false
	//			}
	//
	//			return true
	//		}).Should(BeTrue())
	//
	//		Consistently(func() bool {
	//			obj := &networkingv1.NetworkPolicy{}
	//
	//			objKey := types.NamespacedName{Namespace: "foo2", Name: SampleNetworkPolicyName}
	//			if err := k8sClient.Get(context.TODO(), objKey, obj); err != nil {
	//				return false
	//			}
	//
	//			return true
	//
	//		}).Should(BeFalse())
	//	})
	//
	//	Specify("network policy should allow correspond namespace access", func() {
	//		obj := &networkingv1.NetworkPolicy{}
	//		objKey := types.NamespacedName{Namespace: "foo1", Name: SampleNetworkPolicyName}
	//		err := k8sClient.Get(context.TODO(), objKey, obj)
	//		Expect(err).NotTo(HaveOccurred())
	//		Expect(IsNetworkPolicyAllowFrom(obj, "bar1")).To(BeTrue())
	//		Expect(IsNetworkPolicyAllowFrom(obj, "bar2")).To(BeFalse())
	//	})
	//})
})
