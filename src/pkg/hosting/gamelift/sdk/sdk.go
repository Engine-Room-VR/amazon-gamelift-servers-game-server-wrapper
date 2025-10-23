/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package sdk

import (
	"context"
	"log/slog"

	"github.com/amazon-gamelift/amazon-gamelift-servers-go-server-sdk/model/request"
	"github.com/amazon-gamelift/amazon-gamelift-servers-go-server-sdk/server"
)

// GameLiftSdk defines the interface for interacting with server SDK for Amazon GameLift Servers.
type GameLiftSdk interface {
	// InitSDK initializes the Amazon GameLift SDK.
	InitSDK(ctx context.Context, params server.ServerParameters) error

	// InitSDKFromEnvironment initializes the Amazon GameLift SDK using environment variables for configuration.
	InitSDKFromEnvironment(ctx context.Context) error

	// ProcessReady notifies the Amazon GameLift service that the server process is ready to host game sessions.
	ProcessReady(ctx context.Context, params server.ProcessParameters) error

	// ProcessEnding notifies the Amazon GameLift service that the server process is shutting down.
	ProcessEnding(ctx context.Context) error

	// ActivateGameSession notifies Amazon GameLift that the server process has
	// activated a game session and is now ready to receive player connections.
	ActivateGameSession(ctx context.Context) error

	// GetFleetRoleCredentials retrieves temporary IAM credentials associated with the fleet role.
	GetFleetRoleCredentials(ctx context.Context, roleArn string, roleSessionName string) (accessKeyId string, secretAccessKey string, sessionToken string, err error)

	// Destroy frees the server SDK for Amazon GameLift Servers from memory.
	Destroy(ctx context.Context) error
}

type Sdk struct {
	logger *slog.Logger
}

func (sdk *Sdk) InitSDK(ctx context.Context, params server.ServerParameters) error {
	redactedParams := params
	redactedParams.AuthToken = "<REDACTED>"
	redactedParams.AccessKey = "<REDACTED>"
	redactedParams.SecretKey = "<REDACTED>"
	redactedParams.SessionToken = "<REDACTED>"
	sdk.logger.DebugContext(ctx, "InitSDK called", "params", redactedParams)
	return server.InitSDK(params)
}

func (sdk *Sdk) InitSDKFromEnvironment(ctx context.Context) error {
	sdk.logger.DebugContext(ctx, "InitSDKFromEnvironment called")
	return server.InitSDKFromEnvironment()
}

func (sdk *Sdk) ProcessReady(ctx context.Context, params server.ProcessParameters) error {
	sdk.logger.DebugContext(ctx, "ProcessReady called", "port", params.Port, "logParams", params.LogParameters)
	return server.ProcessReady(params)
}

func (sdk *Sdk) ProcessEnding(ctx context.Context) error {
	sdk.logger.DebugContext(ctx, "ProcessEnding called")
	return server.ProcessEnding()
}

func (sdk *Sdk) ActivateGameSession(ctx context.Context) error {
	sdk.logger.DebugContext(ctx, "ActivateGameSession called")
	return server.ActivateGameSession()
}

func (sdk *Sdk) GetFleetRoleCredentials(ctx context.Context, roleArn string, roleSessionName string) (string, string, string, error) {
	sdk.logger.DebugContext(ctx, "GetFleetRoleCredentials called", "roleArn", roleArn, "roleSessionName", roleSessionName)
	req := request.NewGetFleetRoleCredentials()
	req.RoleArn = roleArn
	req.RoleSessionName = roleSessionName
	creds, err := server.GetFleetRoleCredentials(req)
	if err != nil {
		return "", "", "", err
	}
	return creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken, nil
}

func (sdk *Sdk) Destroy(ctx context.Context) error {
	sdk.logger.DebugContext(ctx, "Destroy called")
	return server.Destroy()
}

func NewSdk(ctx context.Context, logger *slog.Logger) *Sdk {
	s := &Sdk{
		logger: logger,
	}

	server.SetLoggerInterface(NewLogAdaptor(ctx, logger))

	return s
}
