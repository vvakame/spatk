package spatest

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

func getProjectID(projectID string) string {
	if projectID == "" {
		projectID = os.Getenv("SPANNER_EMULATOR_PROJECT_ID")
	}

	return projectID
}

func getInstanceID(instanceID string) string {
	if instanceID == "" {
		instanceID = os.Getenv("SPANNER_EMULATOR_INSTANCE_NAME")
	}

	return instanceID
}

type options struct {
	projectID           string
	instanceID          string
	adminClient         *database.DatabaseAdminClient
	ddls                []string
	preventDropDatabase bool
}

type Option func(opts *options)

func WithProjectID(projectID string) Option {
	return func(opts *options) {
		opts.projectID = projectID
	}
}

func WithInstanceID(instanceID string) Option {
	return func(opts *options) {
		opts.instanceID = instanceID
	}
}

func WithDatabaseAdminClient(adminClient *database.DatabaseAdminClient) Option {
	return func(opts *options) {
		opts.adminClient = adminClient
	}
}

func WithDDLStatements(ddls []string) Option {
	return func(opts *options) {
		opts.ddls = ddls
	}
}

func WithoutDropDatabase() Option {
	return func(opts *options) {
		opts.preventDropDatabase = true
	}
}

func setupDB(ctx context.Context, dbName string, opts *options) (databaseName string, closeFn func() error, err error) {
	var closeFns []func() error
	closeFn = func() error {
		for _, f := range closeFns {
			err := f()
			if err != nil {
				return err
			}
		}
		return nil
	}

	adminClient := opts.adminClient
	if adminClient == nil {
		var err error
		adminClient, err = database.NewDatabaseAdminClient(ctx)
		if err != nil {
			return "", closeFn, err
		}

		f := func() error {
			return adminClient.Close()
		}
		closeFns = append([]func() error{f}, closeFns...)
	}
	projectID := getProjectID(opts.projectID)
	instanceID := getInstanceID(opts.instanceID)
	stmts := opts.ddls

	op, err := adminClient.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID),
		CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", dbName),
		ExtraStatements: stmts,
	})
	if err != nil {
		return "", closeFn, err
	}
	_, err = op.Wait(ctx)
	if err != nil {
		return "", closeFn, err
	}

	databaseName = fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceID, dbName)
	if !opts.preventDropDatabase {
		f := func() error {
			return adminClient.DropDatabase(ctx, &databasepb.DropDatabaseRequest{
				Database: databaseName,
			})
		}
		closeFns = append([]func() error{f}, closeFns...)
	}

	return databaseName, closeFn, nil
}

// SetupDedicatedDB in instance for parallel unit testing.
func SetupDedicatedDB(ctx context.Context, os ...Option) (string, func() error, error) {
	ops := &options{}
	for _, o := range os {
		o(ops)
	}

	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", nil, err
	}

	dbName := fmt.Sprintf("for-ut-%x", b)

	databaseName, closeFn, err := setupDB(ctx, dbName, ops)
	if err != nil {
		return "", nil, err
	}

	return databaseName, closeFn, nil
}
