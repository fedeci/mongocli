// Copyright 2022 MongoDB Inc
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

package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mongodb/mongocli/internal/cli"
	"github.com/mongodb/mongocli/internal/cli/require"
	"github.com/mongodb/mongocli/internal/config"
	"github.com/mongodb/mongocli/internal/flag"
	"github.com/mongodb/mongocli/internal/oauth"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"go.mongodb.org/atlas/auth"
	atlas "go.mongodb.org/atlas/mongodbatlas"
)

//go:generate mockgen -destination=../../mocks/mock_login.go -package=mocks github.com/mongodb/mongocli/internal/cli/auth Authenticator,LoginConfig

type Authenticator interface {
	RequestCode(context.Context) (*auth.DeviceCode, *atlas.Response, error)
	PollToken(context.Context, *auth.DeviceCode) (*auth.Token, *atlas.Response, error)
}

type LoginConfig interface {
	config.SetSaver
	AccessTokenSubject() (string, error)
}

type loginOpts struct {
	cli.DefaultSetterOpts
	AccessToken    string
	RefreshToken   string
	isGov          bool
	isCloudManager bool
	noBrowser      bool
	skipConfig     bool
	config         LoginConfig
	flow           Authenticator
}

func (opts *loginOpts) initFlow() error {
	var err error
	opts.flow, err = oauth.FlowWithConfig(config.Default())
	return err
}

func (opts *loginOpts) SetOAuthUpAccess() {
	switch {
	case opts.isGov:
		opts.Service = config.CloudGovService
	case opts.isCloudManager:
		opts.Service = config.CloudManagerService
	default:
		opts.Service = config.CloudService
	}
	opts.config.Set("service", opts.Service)

	if opts.AccessToken != "" {
		opts.config.Set(config.AccessTokenField, opts.AccessToken)
	}
	if opts.RefreshToken != "" {
		opts.config.Set(config.RefreshTokenField, opts.RefreshToken)
	}
	if opts.OpsManagerURL != "" {
		opts.config.Set(config.OpsManagerURLField, opts.OpsManagerURL)
	}
	if config.ClientID() != "" {
		opts.config.Set(config.ClientIDField, config.ClientID())
	}
}

func (opts *loginOpts) Run(ctx context.Context) error {
	if err := opts.oauthFlow(ctx); err != nil {
		return err
	}
	opts.SetOAuthUpAccess()
	s, err := opts.config.AccessTokenSubject()
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(opts.OutWriter, "Successfully logged in as %s.\n", s)
	if opts.skipConfig {
		return opts.config.Save()
	}
	if err := opts.InitStore(ctx); err != nil {
		return err
	}
	_, _ = fmt.Fprint(opts.OutWriter, "Press Enter to continue your profile configuration")
	_, _ = fmt.Scanln()
	if err := opts.AskOrg(); err != nil {
		return err
	}
	opts.SetUpOrg()
	if err := opts.AskProject(); err != nil {
		return err
	}
	opts.SetUpProject()

	if err := survey.Ask(opts.DefaultQuestions(), opts); err != nil {
		return err
	}
	opts.SetUpOutput()
	opts.SetUpMongoSHPath()
	if err := opts.config.Save(); err != nil {
		return err
	}
	_, _ = fmt.Fprint(opts.OutWriter, "\nYour profile is now configured.\n")
	if config.Name() != config.DefaultProfile {
		_, _ = fmt.Fprintf(opts.OutWriter, "To use this profile, you must set the flag [-%s %s] for every command.\n", flag.ProfileShort, config.Name())
	}

	_, _ = fmt.Fprintf(opts.OutWriter, "You can use [%s config set] to change these settings at a later time.\n", config.BinName())

	return nil
}

func (opts *loginOpts) oauthFlow(ctx context.Context) error {
	code, _, err := opts.flow.RequestCode(ctx)
	if err != nil {
		return err
	}

	codeDuration := time.Duration(code.ExpiresIn) * time.Second
	_, _ = fmt.Fprintf(opts.OutWriter, `
First, copy your one-time code: %s-%s

Next, sign in with your browser and enter the code.

Or go to %s

Your code will expire after %.0f minutes.
`,
		code.UserCode[0:len(code.UserCode)/2],
		code.UserCode[len(code.UserCode)/2:],
		code.VerificationURI,
		codeDuration.Minutes(),
	)
	if !opts.noBrowser {
		if errBrowser := browser.OpenURL(code.VerificationURI); errBrowser != nil {
			_, _ = fmt.Fprintf(os.Stderr, "There was an issue opening your browser\n")
		}
	}

	accessToken, _, err := opts.flow.PollToken(ctx, code)
	if err != nil {
		return err
	}
	opts.AccessToken = accessToken.AccessToken
	opts.RefreshToken = accessToken.RefreshToken
	return nil
}

func hasUserProgrammaticKeys() bool {
	return config.PublicAPIKey() != "" && config.PrivateAPIKey() != ""
}

func LoginBuilder() *cobra.Command {
	opts := &loginOpts{}
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with MongoDB Atlas.",
		Example: fmt.Sprintf(`  To start the interactive setup:
  $ %s auth login
`, config.BinName()),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if hasUserProgrammaticKeys() {
				return fmt.Errorf(`you have already set the programmatic keys for this profile. 

Run '%s auth login --profile <profile_name>' to use your username and password on a new profile`, config.BinName())
			}

			opts.OutWriter = cmd.OutOrStdout()
			opts.config = config.Default()
			if config.OpsManagerURL() != "" {
				opts.OpsManagerURL = config.OpsManagerURL()
			}
			return opts.initFlow()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run(cmd.Context())
		},
		Args: require.NoArgs,
	}

	if config.ToolName == config.MongoCLI {
		cmd.Flags().BoolVar(&opts.isCloudManager, "cm", false, "Log in to Cloud Manager.")
	}

	cmd.Flags().BoolVar(&opts.isGov, "gov", false, "Log in to Atlas for Government.")
	cmd.Flags().BoolVar(&opts.noBrowser, "noBrowser", false, "Don't try to open a browser session.")
	cmd.Flags().BoolVar(&opts.skipConfig, "skipConfig", false, "Skip profile configuration.")
	return cmd
}

func Builder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage the CLI's authentication state.",
		Annotations: map[string]string{
			"toc": "true",
		},
	}
	cmd.AddCommand(
		LoginBuilder(),
		WhoAmIBuilder(),
		LogoutBuilder(),
	)

	return cmd
}
