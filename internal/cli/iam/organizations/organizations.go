// Copyright 2020 MongoDB Inc
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

package organizations

import (
	"github.com/mongodb/mongocli/internal/cli"
	"github.com/mongodb/mongocli/internal/cli/iam/organizations/apikeys"
	"github.com/mongodb/mongocli/internal/cli/iam/organizations/invitations"
	"github.com/mongodb/mongocli/internal/cli/iam/organizations/users"
	"github.com/spf13/cobra"
)

func Builder() *cobra.Command {
	const use = "organizations"
	cmd := &cobra.Command{
		Use:     use,
		Short:   "Organization operations.",
		Long:    "Create, list and manage your MongoDB organizations.",
		Aliases: cli.GenerateAliases(use, "orgs", "org"),
	}
	cmd.AddCommand(
		ListBuilder(),
		DescribeBuilder(),
		CreateBuilder(),
		DeleteBuilder(),
		apikeys.Builder(),
		users.Builder(),
		invitations.Builder(),
	)
	return cmd
}

func AtlasCLIBuilder() *cobra.Command {
	const use = "organizations"
	cmd := &cobra.Command{
		Use:     use,
		Short:   "Organization operations.",
		Long:    "Create, list and manage your MongoDB organizations.",
		Aliases: cli.GenerateAliases(use, "orgs", "org"),
	}
	cmd.AddCommand(
		ListBuilder(),
		DescribeBuilder(),
		DeleteBuilder(),
		apikeys.Builder(),
		users.Builder(),
		invitations.Builder(),
	)
	return cmd
}
