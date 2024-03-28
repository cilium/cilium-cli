// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package encrypt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/cilium/cilium/api/v1/models"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cilium/cilium-cli/defaults"
	"github.com/cilium/cilium-cli/status"
	"github.com/cilium/cilium-cli/utils/features"
)

// GetEncryptStatus gets encryption status from all/specific cilium agent pods.
func (s *Encrypt) GetEncryptStatus(ctx context.Context) error {
	ctx, cancelFn := context.WithTimeout(ctx, s.params.WaitDuration)
	defer cancelFn()

	pods, err := s.fetchCiliumPods(ctx)
	if err != nil {
		return err
	}

	statusMap, err := s.fetchEncryptStatusConcurrently(ctx, pods)
	if err != nil {
		return err
	}

	result, err := s.enrichWithIPsecDetails(ctx, statusMap)
	if err != nil {
		return err
	}

	return s.writeStatus(result)
}

func (s *Encrypt) fetchEncryptStatusConcurrently(ctx context.Context, pods []corev1.Pod) (map[string]models.EncryptionStatus, error) {
	// res contains data returned from cilium pod
	type res struct {
		nodeName string
		status   models.EncryptionStatus
		err      error
	}
	resCh := make(chan res)
	defer close(resCh)

	// concurrently fetch state from each cilium pod
	for _, pod := range pods {
		go func(ctx context.Context, pod corev1.Pod) {
			st, err := s.fetchEncryptStatusFromPod(ctx, pod)
			resCh <- res{
				nodeName: pod.Spec.NodeName,
				status:   st,
				err:      err,
			}
		}(ctx, pod)
	}

	// read from the channel, on error, store error and continue to next node
	var err error
	data := make(map[string]models.EncryptionStatus)
	for range pods {
		r := <-resCh
		if r.err != nil {
			err = errors.Join(err, r.err)
			continue
		}
		data[r.nodeName] = r.status
	}
	return data, err
}

func (s *Encrypt) fetchEncryptStatusFromPod(ctx context.Context, pod corev1.Pod) (models.EncryptionStatus, error) {
	cmd := []string{"cilium", "encrypt", "status", "-o", "json"}
	output, err := s.client.ExecInPod(ctx, pod.Namespace, pod.Name, defaults.AgentContainerName, cmd)
	if err != nil {
		return models.EncryptionStatus{}, fmt.Errorf("failed to fetch encryption status from %s: %v", pod.Name, err)
	}
	encStatus, err := nodeStatusFromOutput(output.String())
	if err != nil {
		return models.EncryptionStatus{}, fmt.Errorf("failed to parse encryption status from %s: %v", pod.Name, err)
	}
	return encStatus, nil
}

func nodeStatusFromOutput(output string) (models.EncryptionStatus, error) {
	if !json.Valid([]byte(output)) {
		res, err := nodeStatusFromText(output)
		if err != nil {
			return models.EncryptionStatus{}, fmt.Errorf("failed to parse text: %v", err)
		}
		return res, nil
	}
	encStatus := models.EncryptionStatus{}
	if err := json.Unmarshal([]byte(output), &encStatus); err != nil {
		return models.EncryptionStatus{}, fmt.Errorf("failed to unmarshal json: %v", err)
	}
	return encStatus, nil
}

func nodeStatusFromText(str string) (models.EncryptionStatus, error) {
	res := models.EncryptionStatus{
		Ipsec: &models.IPsecStatus{
			DecryptInterfaces: make([]string, 0),
			XfrmErrors:        make(map[string]int64),
		},
		Wireguard: &models.WireguardStatus{
			Interfaces: make([]*models.WireguardInterface, 0),
		},
	}
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "Encryption":
			res.Mode = value
		case "Max Seq. Number":
			res.Ipsec.MaxSeqNumber = value
		case "Decryption interface(s)":
			if value != "" {
				res.Ipsec.DecryptInterfaces = append(res.Ipsec.DecryptInterfaces, value)
			}
		case "Keys in use":
			keys, err := strconv.Atoi(value)
			if err != nil {
				return models.EncryptionStatus{}, fmt.Errorf("invalid number 'Keys in use' [%s]: %v", value, err)
			}
			res.Ipsec.KeysInUse = int64(keys)
		case "Errors":
			count, err := strconv.Atoi(value)
			if err != nil {
				return models.EncryptionStatus{}, fmt.Errorf("invalid number 'Errors' [%s]: %v", value, err)
			}
			res.Ipsec.ErrorCount = int64(count)
		default:
			count, err := strconv.Atoi(value)
			if err != nil {
				return models.EncryptionStatus{}, fmt.Errorf("invalid number '%s' [%s]: %v", key, value, err)
			}
			res.Ipsec.XfrmErrors[key] = int64(count)
		}
	}
	return res, nil
}

func (s *Encrypt) enrichWithIPsecDetails(ctx context.Context, statusMap map[string]models.EncryptionStatus) (map[string]encryptionStatus, error) {
	ipsecEnabled := false
	result := make(map[string]encryptionStatus, len(statusMap))
	for k, s := range statusMap {
		if s.Ipsec != nil {
			ipsecEnabled = true
		}
		result[k] = encryptionStatus{EncryptionStatus: s}
	}

	if !ipsecEnabled {
		return result, nil
	}

	cm, err := s.client.GetConfigMap(ctx, s.params.CiliumNamespace, defaults.ConfigMapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable get ConfigMap %q: %w", defaults.ConfigMapName, err)
	}
	fs := features.Set{}
	fs.ExtractFromConfigMap(cm)

	key, err := s.readIPsecKey(ctx)
	if err != nil {
		return nil, err
	}
	perNodeIPsecKey := strings.Contains(key, "+")
	expectedKeyCount := expectedIPsecKeyCount(len(statusMap), fs, perNodeIPsecKey)

	for k, s := range result {
		if s.Ipsec == nil {
			continue
		}
		s.IPsecExpectedKeyCount = expectedKeyCount
		s.IPsecPerNodeKey = perNodeIPsecKey
		s.IPsecKeyRotationInProgress = int64(expectedKeyCount) != s.Ipsec.KeysInUse
		result[k] = s
	}
	return result, nil
}

func expectedIPsecKeyCount(ciliumPods int, fs features.Set, perNodeKey bool) int {
	if !perNodeKey {
		return 1
	}
	// IPsec key states for `local_cilium_internal_ip` and `remote_node_ip`
	xfrmStates := 2
	if fs[features.UseCiliumInternalIPForIPsec].Enabled {
		// plus state for `remote_cilium_internal_ip`
		xfrmStates++
	}
	if fs[features.IPv6].Enabled {
		// multiply by 2 because of dual state: IPv4 & IPv6
		xfrmStates *= 2
	}
	// subtract 1 to count remote nodes only
	return (ciliumPods - 1) * xfrmStates
}

func (s *Encrypt) writeStatus(res map[string]encryptionStatus) error {
	if s.params.PerNodeDetails {
		for nodeName, n := range res {
			if err := printStatus(nodeName, n, s.params.Output); err != nil {
				return err
			}
		}
		return nil
	}
	cs, err := clusterNodeStatus(res)
	if err != nil {
		return err
	}
	return cs.printStatus(s.params.Output)
}

func clusterNodeStatus(res map[string]encryptionStatus) (clusterStatus, error) {
	cs := clusterStatus{
		TotalNodeCount:          len(res),
		IPsecKeysInUseNodeCount: make(map[int64]int64),
		XfrmErrors:              make(map[string]int64),
		XfrmErrorNodeCount:      make(map[string]int64),
	}
	for _, v := range res {
		if v.Mode == "Disabled" {
			cs.EncDisabledNodeCount++
			continue
		}
		if v.Mode == "Wireguard" {
			cs.EncWireguardNodeCount++
			continue
		}
		if v.Mode == "IPsec" {
			cs.EncIPsecNodeCount++
			cs.IPsecExpectedKeyCount = v.IPsecExpectedKeyCount
			cs.IPsecPerNodeKey = v.IPsecPerNodeKey
			cs.IPsecKeyRotationInProgress = int64(v.IPsecExpectedKeyCount) != v.Ipsec.KeysInUse
		}
		cs.IPsecKeysInUseNodeCount[v.Ipsec.KeysInUse]++
		maxSeqNum, err := maxSequenceNumber(v.Ipsec.MaxSeqNumber, cs.IPsecMaxSeqNum)
		if err != nil {
			return clusterStatus{}, err
		}
		cs.IPsecMaxSeqNum = maxSeqNum
		cs.IPsecErrCount += v.Ipsec.ErrorCount
		for k, e := range v.Ipsec.XfrmErrors {
			cs.XfrmErrors[k] += e
			cs.XfrmErrorNodeCount[k]++
		}
	}
	return cs, nil
}

func (c clusterStatus) printStatus(format string) error {
	if format == status.OutputJSON {
		return printJSONStatus(c)
	}

	builder := strings.Builder{}
	if c.EncDisabledNodeCount > 0 {
		builder.WriteString(fmt.Sprintf("Encryption: Disabled (%d/%d nodes)\n", c.EncDisabledNodeCount, c.TotalNodeCount))
	}
	if c.EncIPsecNodeCount > 0 {
		builder.WriteString(fmt.Sprintf("Encryption: IPsec (%d/%d nodes)\n", c.EncIPsecNodeCount, c.TotalNodeCount))
		if len(c.IPsecKeysInUseNodeCount) > 0 {
			keys := make([]int64, 0, len(c.IPsecKeysInUseNodeCount))
			for k := range c.IPsecKeysInUseNodeCount {
				keys = append(keys, k)
			}
			sort.Slice(keys, func(i, j int) bool {
				return keys[i] < keys[j]
			})
			keyStrs := make([]string, 0, len(keys))
			for _, k := range keys {
				keyStrs = append(keyStrs, fmt.Sprintf("%d on %d/%d", k, c.IPsecKeysInUseNodeCount[k], c.TotalNodeCount))
			}
			builder.WriteString(fmt.Sprintf("IPsec keys in use: %s\n", strings.Join(keyStrs, ", ")))
		}
		builder.WriteString(fmt.Sprintf("IPsec highest Seq. Number: %s across all nodes\n", c.IPsecMaxSeqNum))
		builder.WriteString(fmt.Sprintf("IPsec expected key count: %d\n", c.IPsecExpectedKeyCount))
		builder.WriteString(fmt.Sprintf("IPsec key rotation in progress: %t\n", c.IPsecKeyRotationInProgress))
		builder.WriteString(fmt.Sprintf("IPsec per-node key: %t\n", c.IPsecPerNodeKey))
		builder.WriteString(fmt.Sprintf("IPsec errors: %d across all nodes\n", c.IPsecErrCount))
		for k, v := range c.XfrmErrors {
			builder.WriteString(fmt.Sprintf("\t%s: %d on %d/%d nodes\n", k, v, c.XfrmErrorNodeCount[k], c.TotalNodeCount))
		}
	}
	if c.EncWireguardNodeCount > 0 {
		builder.WriteString(fmt.Sprintf("Encryption: Wireguard (%d/%d nodes)\n", c.EncWireguardNodeCount, c.TotalNodeCount))
	}
	_, err := fmt.Println(builder.String())
	return err
}

func printStatus(nodeName string, n encryptionStatus, format string) error {
	if format == status.OutputJSON {
		return printJSONStatus(n)
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Node: %s\n", nodeName))
	builder.WriteString(fmt.Sprintf("Encryption: %s\n", n.Mode))
	if n.Mode == "IPsec" {
		if n.Ipsec.KeysInUse > 0 {
			builder.WriteString(fmt.Sprintf("IPsec keys in use: %d\n", n.Ipsec.KeysInUse))
		}
		if n.Ipsec.MaxSeqNumber != "" {
			builder.WriteString(fmt.Sprintf("IPsec highest Seq. Number: %s\n", n.Ipsec.MaxSeqNumber))
		}
		builder.WriteString(fmt.Sprintf("IPsec expected key count: %d\n", n.IPsecExpectedKeyCount))
		builder.WriteString(fmt.Sprintf("IPsec key rotation in progress: %t\n", n.IPsecKeyRotationInProgress))
		builder.WriteString(fmt.Sprintf("IPsec per-node key: %t\n", n.IPsecPerNodeKey))
		builder.WriteString(fmt.Sprintf("IPsec errors: %d\n", n.Ipsec.ErrorCount))
		for k, v := range n.Ipsec.XfrmErrors {
			builder.WriteString(fmt.Sprintf("\t%s: %d\n", k, v))
		}
	}
	_, err := fmt.Println(builder.String())
	return err
}

func printJSONStatus(v any) error {
	js, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}
	fmt.Println(string(js))
	return nil
}
