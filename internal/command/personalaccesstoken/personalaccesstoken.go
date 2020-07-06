package personalaccesstoken

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.com/gitlab-org/gitlab-shell/internal/command/commandargs"
	"gitlab.com/gitlab-org/gitlab-shell/internal/command/readwriter"
	"gitlab.com/gitlab-org/gitlab-shell/internal/config"
	"gitlab.com/gitlab-org/gitlab-shell/internal/gitlabnet/personalaccesstoken"
)

const (
	usageText         = "Usage: personal_access_token <name> <scope1[,scope2,...]> [ttl_days]"
	expiresDateFormat = "2006-01-02"
)

type Command struct {
	Config     *config.Config
	Args       *commandargs.Shell
	ReadWriter *readwriter.ReadWriter
	TokenArgs  *tokenArgs
}

type tokenArgs struct {
	Name        string
	Scopes      []string
	ExpiresDate string // Calculated, a TTL is passed from command-line.
}

func (c *Command) Execute() error {
	err := c.parseTokenArgs()
	if err != nil {
		return err
	}

	response, err := c.getPersonalAccessToken()
	if err != nil {
		return err
	}

	fmt.Fprint(c.ReadWriter.Out, "Token:   "+response.Token+"\n")
	fmt.Fprint(c.ReadWriter.Out, "Scopes:  "+strings.Join(response.Scopes, ",")+"\n")
	if response.ExpiresAt == "" {
		fmt.Fprint(c.ReadWriter.Out, "Expires: never\n")
	} else {
		fmt.Fprint(c.ReadWriter.Out, "Expires: "+response.ExpiresAt+"\n")
	}
	return nil
}

func (c *Command) parseTokenArgs() error {
	if len(c.Args.SshArgs) < 3 || len(c.Args.SshArgs) > 4 {
		return errors.New(usageText)
	}
	c.TokenArgs = &tokenArgs{
		Name:   c.Args.SshArgs[1],
		Scopes: strings.Split(c.Args.SshArgs[2], ","),
	}

	if len(c.Args.SshArgs) < 4 {
		return nil
	}
	rawTTL := c.Args.SshArgs[3]

	TTL, err := strconv.Atoi(rawTTL)
	if err != nil || TTL < 0 {
		return fmt.Errorf("Invalid value for days_ttl: '%s'", rawTTL)
	}

	c.TokenArgs.ExpiresDate = time.Now().AddDate(0, 0, TTL+1).Format(expiresDateFormat)

	return nil
}

func (c *Command) getPersonalAccessToken() (*personalaccesstoken.Response, error) {
	client, err := personalaccesstoken.NewClient(c.Config)
	if err != nil {
		return nil, err
	}

	return client.GetPersonalAccessToken(c.Args, c.TokenArgs.Name, &c.TokenArgs.Scopes, c.TokenArgs.ExpiresDate)
}
