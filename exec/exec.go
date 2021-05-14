// Copyright 2021 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exec

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ciliumExecImplementation interface {
	DiscoverCiliumNamespace(ctx context.Context) (string, error)
	ExecInPod(ctx context.Context, namespace, name, container string, command []string, interactive bool) (bytes.Buffer, bytes.Buffer, error)
	GetNode(ctx context.Context, name string, opts metav1.GetOptions) (*corev1.Node, error)
	GetPod(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Pod, error)
	ListPods(ctx context.Context, namespace string, opts metav1.ListOptions) (*corev1.PodList, error)
}

type CiliumExec struct {
	client ciliumExecImplementation
	params Parameters
}

type Parameters struct {
	Target  string
	Command []string
}

func NewCiliumExecImplementation(client ciliumExecImplementation, p Parameters) (*CiliumExec, error) {
	return &CiliumExec{
		client: client,
		params: p,
	}, nil
}

// ExecInTarget attempts to find the appropriate target and executes the specified command.
func (ce *CiliumExec) ExecInTarget(ctx context.Context) error {
	// Start by attempting to discover the namespace in which Cilium is installed.
	ns, err := ce.client.DiscoverCiliumNamespace(ctx)
	if err != nil {
		return err
	}
	// Try to understand if the specified target string is...
	// 1. ... the name of a node;
	// 2. ... the name of a Cilium pod;
	// 3. ... a '<namespace>/<name>' key targeting a pod.
	// Pods running in the 'default' namespace need to be specified as 'default/<pod-name>'.
	pp := strings.Split(ce.params.Target, "/")
	if err != nil {
		return fmt.Errorf("failed to parse %q as a target: %e", ce.params.Target, err)
	}
	var pn string
	switch len(pp) {
	case 1:
		// Check whether a node with the provided name exists and use that.
		e, err := ce.nodeExists(ctx, pp[0])
		if err != nil {
			return err
		}
		// If no node with the specified name exists, we assume the target is the name of a Cilium pod.
		if !e {
			pn = pp[0]
			break
		}
		// Lookup the name of the Cilium pod running on the node.
		p, err := ce.getCiliumPodInNode(ctx, ns, pp[0])
		if err != nil {
			return err
		}
		pn = p
	case 2:
		// Lookup the name of the node where the pod referenced by the target string is running.
		cn, err := ce.getNodeNameFromPod(ctx, pp[0], pp[1])
		if err != nil {
			return fmt.Errorf("failed to infer target node name: %w", err)
		}
		// Lookup the name of the Cilium pod running on the node.
		p, err := ce.getCiliumPodInNode(ctx, ns, cn)
		if err != nil {
			return err
		}
		pn = p
	default:
		return fmt.Errorf("%q is not a valid target format", ce.params.Target)
	}

	// Double-check whether the targeted pod exists and is running.
	p, err := ce.client.GetPod(ctx, ns, pn, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return err
		}
		return fmt.Errorf("failed to get Cilium pod %q: %w", pn, err)
	}
	_, _, err = ce.client.ExecInPod(ctx, ns, p.Name, "cilium-agent", ce.params.Command, true)
	return err
}

// getCiliumPodInNode returns the name of the Cilium pod running in the specified node.
func (ce *CiliumExec) getCiliumPodInNode(ctx context.Context, ciliumNamespace, nodeName string) (string, error) {
	p, err := ce.client.ListPods(ctx, ciliumNamespace, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName==%s", nodeName),
		LabelSelector: "k8s-app=cilium",
		Limit:         1,
	})
	if err != nil {
		return "", fmt.Errorf("failed to discover Cilium pod in node %q: %w", nodeName, err)
	}
	if len(p.Items) != 1 {
		return "", fmt.Errorf("found %d Cilium pods running on node %q", len(p.Items), nodeName)
	}
	return p.Items[0].Name, nil
}

// getNodeNameFromPod returns the name of the node where a given pod is running.
func (ce *CiliumExec) getNodeNameFromPod(ctx context.Context, namespace, name string) (string, error) {
	p, err := ce.client.GetPod(ctx, namespace, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return "", fmt.Errorf("no pod named %q exists in namespace %q", name, namespace)
		}
		return "", fmt.Errorf("failed to get pod %q in namespace %q: %w", name, namespace, err)
	}
	if p.Spec.NodeName == "" {
		return "", fmt.Errorf("pod %q in namespace %q hasn't been assigned to a node yet", name, namespace)
	}
	return p.Spec.NodeName, nil
}

// nodeExists returns whether a node with the specified name exists.
func (ce *CiliumExec) nodeExists(ctx context.Context, name string) (bool, error) {
	_, err := ce.client.GetNode(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get node %q: %w", name, err)
	}
	return true, nil
}
