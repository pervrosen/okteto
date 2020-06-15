// Copyright 2020 The Okteto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/okteto/okteto/cmd/utils"
	"github.com/okteto/okteto/pkg/k8s/pods"
	"github.com/okteto/okteto/pkg/log"

	k8Client "github.com/okteto/okteto/pkg/k8s/client"
	"github.com/okteto/okteto/pkg/model"

	"github.com/spf13/cobra"
)

//Restart restarts the pods of a given dev mode deployment
func Restart() *cobra.Command {
	var namespace string
	var devPath string

	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Restarts the deployments listed in the services field of the okteto manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			dev, err := utils.LoadDev(devPath)
			if err != nil {
				return err
			}
			if err := dev.UpdateNamespace(namespace); err != nil {
				return err
			}
			serviceName := ""
			if len(args) > 0 {
				serviceName = args[0]
			}
			if err := executeRestart(dev, serviceName); err != nil {
				return err
			}
			log.Success("Deployments restarted")

			return nil
		},
	}

	cmd.Flags().StringVarP(&devPath, "file", "f", defaultManifest, "path to the manifest file")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace where the restart command is executed")

	return cmd
}

func executeRestart(dev *model.Dev, sn string) error {
	log.Infof("restarting services")
	client, _, namespace, err := k8Client.GetLocal()
	if err != nil {
		return err
	}

	if dev.Namespace == "" {
		dev.Namespace = namespace
	}

	spinner := utils.NewSpinner("Restarting deployments...")
	spinner.Start()
	defer spinner.Stop()

	if err := pods.Restart(dev, client, sn); err != nil {
		return err
	}

	return nil
}
