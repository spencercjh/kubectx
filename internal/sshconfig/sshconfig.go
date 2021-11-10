// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sshconfig

import (
	"bufio"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"sshctx/internal/printer"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Loader interface {
	LoadSSHConfig() (io.ReadWriteCloser, error)
	LoadSSHCTXData() (io.ReadWriteCloser, error)
}

type SSHConfig struct {
	loader        Loader
	sshctxDataRWC io.ReadWriteCloser
	sshconfigRWC  io.ReadWriteCloser
	Hosts         []Host
	PreviousHost  Host
	rootNode      *yaml.Node
}

type Host struct {
	Host     string
	Username string
	Port     int
}

func (h *Host) ToSSHParameter() string {
	return h.Username + "@" + h.Host + ":" + strconv.Itoa(h.Port)
}

var EmptyHost = Host{}

type SSHCTXData struct {
	previous Host
}

func (s *SSHConfig) WithLoader(l Loader) *SSHConfig {
	s.loader = l
	return s
}

func (s *SSHConfig) Close() []error {
	switch {
	case s.sshconfigRWC == nil && s.sshctxDataRWC != nil:
		return []error{nil, s.sshctxDataRWC.Close()}
	case s.sshctxDataRWC == nil && s.sshconfigRWC != nil:
		return []error{s.sshconfigRWC.Close(), nil}
	case s.sshconfigRWC == nil && s.sshctxDataRWC == nil:
		return nil
	}
	return []error{s.sshconfigRWC.Close(), s.sshctxDataRWC.Close()}
}

func (s *SSHConfig) Parse() error {
	if s.loader == nil {
		return errors.New("Missing loader")
	}
	sshconfig, err := s.loader.LoadSSHConfig()
	if err != nil {
		return errors.Wrap(err, "failed to load sshconfig")
	}

	s.sshconfigRWC = sshconfig

	sshctxData, err := s.loader.LoadSSHCTXData()
	if err != nil {
		return errors.Wrap(err, "failed to load sshconfig")
	}
	s.sshctxDataRWC = sshctxData

	s.Hosts, err = getSSHConfigItems(s.sshconfigRWC)
	if err != nil {
		return errors.Wrap(err, "Can not parse sshconfig")
	}

	var v yaml.Node
	if err := yaml.NewDecoder(s.sshctxDataRWC).Decode(&v); err != nil {
		_ = printer.Warning(os.Stderr, "failed to decode sshctxData because :%v", err)
	}
	if len(v.Content) == 0 {
		_ = printer.Warning(os.Stderr, "No previous config")
		s.PreviousHost = EmptyHost
	} else {
		s.rootNode = v.Content[0]
		if s.rootNode.Kind != yaml.MappingNode {
			return errors.New("sshctxData file is not a map document")
		}
		h, err := previousConfig(s.rootNode)
		if err != nil {
			_ = printer.Warning(os.Stderr, "Fail to load previous config: %v", err)
		}
		if h == EmptyHost {
			_ = printer.Warning(os.Stderr, "No previous config")
		}
		s.PreviousHost = h
	}
	return nil
}

func getSSHConfigItems(rwc io.Reader) ([]Host, error) {
	s := bufio.NewScanner(rwc)

	rows, hostIndices, _ := scanSSHConfig(s)
	if len(rows) == 0 || len(hostIndices) == 0 {
		return nil, errors.New("No host found in sshconfig")
	}

	var hosts []Host
	j := 1
	for i := 0; i < len(rows); j++ {
		itemBeginIndex := i
		var itemEndIndex int
		if j >= len(hostIndices) {
			itemEndIndex = len(rows) - 1
		} else {
			itemEndIndex = hostIndices[j]
		}
		if itemBeginIndex == itemEndIndex {
			break
		}
		configItem := extractConfigItem(itemBeginIndex, itemEndIndex, rows)
		if configItem != EmptyHost {
			hosts = append(hosts, configItem)
		}
		i = itemEndIndex
	}
	return hosts, nil
}

func extractConfigItem(itemBeginIndex int, itemEndIndex int, rows []string) Host {
	configItem := Host{}
	var host string
	var hostname string
	for k := itemBeginIndex; k < itemEndIndex; k++ {
		itemRow := rows[k]
		if strings.HasPrefix(itemRow, "Host") {
			host = strings.TrimSpace(itemRow[4:])
		}
		if strings.HasPrefix(itemRow, "Hostname") {
			hostname = strings.TrimSpace(itemRow[8:])
		}
		if strings.HasPrefix(itemRow, "User") {
			configItem.Username = strings.TrimSpace(itemRow[4:])
		}
		if strings.HasPrefix(itemRow, "Port") {
			configItem.Port, _ = strconv.Atoi(strings.TrimSpace(itemRow[4:]))
		}
	}
	if hostname != "" {
		configItem.Host = hostname
	} else if host != "" {
		configItem.Host = host
	}
	// no host and hostname
	if configItem.Host == "name" ||
		// no host and hostname
		configItem.Host == "" ||
		// `*` isn't a host
		configItem.Host == "*" {
		return EmptyHost
	}
	if configItem.Username == "" {
		configItem.Username = os.Getenv("USER")
	}
	if configItem.Port == 0 {
		configItem.Port = 22
	}
	return configItem
}

func scanSSHConfig(scanner *bufio.Scanner) ([]string, []int, error) {
	var rows []string
	var hostIndices []int
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		rows = append(rows, strings.TrimSpace(line))
		if strings.HasPrefix(line, "Host") {
			hostIndices = append(hostIndices, i)
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, errors.Wrap(err, "Can not scan sshconfig")
	}

	return rows, hostIndices, nil
}

func previousConfig(rootNode *yaml.Node) (Host, error) {
	previous := valueOf(rootNode, "previous")
	if previous == nil {
		return EmptyHost, errors.New("\"previous\" entry is nil")
	} else if previous.Kind != yaml.MappingNode {
		return EmptyHost, errors.New("\"previous\" is not a scalar node")
	}
	host := Host{}
	host.Host = valueOf(previous, "host").Value
	host.Username = valueOf(previous, "username").Value
	port, err := strconv.Atoi(valueOf(previous, "port").Value)
	if err != nil {
		return EmptyHost, errors.Wrap(err, "Can't parse port in the previous node")
	}
	host.Port = port
	return host, nil
}

func valueOf(mapNode *yaml.Node, key string) *yaml.Node {
	if mapNode.Kind != yaml.MappingNode {
		return nil
	}
	for i, ch := range mapNode.Content {
		if i%2 == 0 && ch.Kind == yaml.ScalarNode && ch.Value == key {
			return mapNode.Content[i+1]
		}
	}
	return nil
}
