package godo

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	databaseBasePath        = "/v2/databases"
	databaseSinglePath      = databaseBasePath + "/%s"
	databaseResizePath      = databaseBasePath + "/%s/resize"
	databaseMigratePath     = databaseBasePath + "/%s/migrate"
	databaseMaintenancePath = databaseBasePath + "/%s/maintenance"
	databaseBackupsPath     = databaseBasePath + "/%s/backups"
	databaseUsersPath       = databaseBasePath + "/%s/users"
	databaseUserPath        = databaseBasePath + "/%s/users/%s"
	databaseDBPath          = databaseBasePath + "/%s/dbs/%s"
	databaseDBsPath         = databaseBasePath + "/%s/dbs"
	databasePoolPath        = databaseBasePath + "/%s/pools/%s"
	databasePoolsPath       = databaseBasePath + "/%s/pools"
	databaseReplicaPath     = databaseBasePath + "/%s/replicas/%s"
	databaseReplicasPath    = databaseBasePath + "/%s/replicas"
	evictionPolicyPath      = databaseBasePath + "/%s/eviction_policy"
)

// DatabasesService is an interface for interfacing with the databases endpoints
// of the DigitalOcean API.
// See: https://developers.digitalocean.com/documentation/v2#databases
type DatabasesService interface {
	List(context.Context, *ListOptions) ([]Database, *Response, error)
	Get(context.Context, string) (*Database, *Response, error)
	Create(context.Context, *DatabaseCreateRequest) (*Database, *Response, error)
	Delete(context.Context, string) (*Response, error)
	Resize(context.Context, string, *DatabaseResizeRequest) (*Response, error)
	Migrate(context.Context, string, *DatabaseMigrateRequest) (*Response, error)
	UpdateMaintenance(context.Context, string, *DatabaseUpdateMaintenanceRequest) (*Response, error)
	ListBackups(context.Context, string, *ListOptions) ([]DatabaseBackup, *Response, error)
	GetUser(context.Context, string, string) (*DatabaseUser, *Response, error)
	ListUsers(context.Context, string, *ListOptions) ([]DatabaseUser, *Response, error)
	CreateUser(context.Context, string, *DatabaseCreateUserRequest) (*DatabaseUser, *Response, error)
	DeleteUser(context.Context, string, string) (*Response, error)
	ListDBs(context.Context, string, *ListOptions) ([]DatabaseDB, *Response, error)
	CreateDB(context.Context, string, *DatabaseCreateDBRequest) (*DatabaseDB, *Response, error)
	GetDB(context.Context, string, string) (*DatabaseDB, *Response, error)
	DeleteDB(context.Context, string, string) (*Response, error)
	ListPools(context.Context, string, *ListOptions) ([]DatabasePool, *Response, error)
	CreatePool(context.Context, string, *DatabaseCreatePoolRequest) (*DatabasePool, *Response, error)
	GetPool(context.Context, string, string) (*DatabasePool, *Response, error)
	DeletePool(context.Context, string, string) (*Response, error)
	GetReplica(context.Context, string, string) (*DatabaseReplica, *Response, error)
	ListReplicas(context.Context, string, *ListOptions) ([]DatabaseReplica, *Response, error)
	CreateReplica(context.Context, string, *DatabaseCreateReplicaRequest) (*DatabaseReplica, *Response, error)
	DeleteReplica(context.Context, string, string) (*Response, error)
	GetEvictionPolicy(context.Context, string) (string, *Response, error)
	SetEvictionPolicy(context.Context, string, string) (*Response, error)
}

// DatabasesServiceOp handles communication with the Databases related methods
// of the DigitalOcean API.
type DatabasesServiceOp struct {
	client *Client
}

var _ DatabasesService = &DatabasesServiceOp{}

// Database represents a DigitalOcean managed database product. These managed databases
// are usually comprised of a cluster of database nodes, a primary and 0 or more replicas.
// The EngineSlug is a string which indicates the type of database service. Some examples are
// "pg", "mysql" or "redis". A Database also includes connection information and other
// properties of the service like region, size and current status.
type Database struct {
	ID                 string                     `json:"id,omitempty"`
	Name               string                     `json:"name,omitempty"`
	EngineSlug         string                     `json:"engine,omitempty"`
	VersionSlug        string                     `json:"version,omitempty"`
	Connection         *DatabaseConnection        `json:"connection,omitempty"`
	PrivateConnection  *DatabaseConnection        `json:"private_connection,omitempty"`
	Users              []DatabaseUser             `json:"users,omitempty"`
	NumNodes           int                        `json:"num_nodes,omitempty"`
	SizeSlug           string                     `json:"size,omitempty"`
	DBNames            []string                   `json:"db_names,omitempty"`
	RegionSlug         string                     `json:"region,omitempty"`
	Status             string                     `json:"status,omitempty"`
	MaintenanceWindow  *DatabaseMaintenanceWindow `json:"maintenance_window,omitempty"`
	CreatedAt          time.Time                  `json:"created_at,omitempty"`
	PrivateNetworkUUID string                     `json:"private_network_uuid,omitempty"`
}

// DatabaseConnection represents a database connection
type DatabaseConnection struct {
	URI      string `json:"uri,omitempty"`
	Database string `json:"database,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	SSL      bool   `json:"ssl,omitempty"`
}

// DatabaseUser represents a user in the database
type DatabaseUser struct {
	Name     string `json:"name,omitempty"`
	Role     string `json:"role,omitempty"`
	Password string `json:"password,omitempty"`
}

// DatabaseMaintenanceWindow represents the maintenance_window of a database
// cluster
type DatabaseMaintenanceWindow struct {
	Day         string   `json:"day,omitempty"`
	Hour        string   `json:"hour,omitempty"`
	Pending     bool     `json:"pending,omitempty"`
	Description []string `json:"description,omitempty"`
}

// DatabaseBackup represents a database backup.
type DatabaseBackup struct {
	CreatedAt     time.Time `json:"created_at,omitempty"`
	SizeGigabytes float64   `json:"size_gigabytes,omitempty"`
}

// DatabaseCreateRequest represents a request to create a database cluster
type DatabaseCreateRequest struct {
	Name       string `json:"name,omitempty"`
	EngineSlug string `json:"engine,omitempty"`
	Version    string `json:"version,omitempty"`
	SizeSlug   string `json:"size,omitempty"`
	Region     string `json:"region,omitempty"`
	NumNodes   int    `json:"num_nodes,omitempty"`
}

// DatabaseResizeRequest can be used to initiate a database resize operation.
type DatabaseResizeRequest struct {
	SizeSlug string `json:"size,omitempty"`
	NumNodes int    `json:"num_nodes,omitempty"`
}

// DatabaseMigrateRequest can be used to initiate a database migrate operation.
type DatabaseMigrateRequest struct {
	Region string `json:"region,omitempty"`
}

// DatabaseUpdateMaintenanceRequest can be used to update the database's maintenance window.
type DatabaseUpdateMaintenanceRequest struct {
	Day  string `json:"day,omitempty"`
	Hour string `json:"hour,omitempty"`
}

// DatabaseDB represents an engine-specific database created within a database cluster. For SQL
// databases like PostgreSQL or MySQL, a "DB" refers to a database created on the RDBMS. For instance,
// a PostgreSQL database server can contain many database schemas, each with it's own settings, access
// permissions and data. ListDBs will return all databases present on the server.
type DatabaseDB struct {
	Name string `json:"name"`
}

// DatabaseReplica represents a read-only replica of a particular database
type DatabaseReplica struct {
	Name               string              `json:"name"`
	Connection         *DatabaseConnection `json:"connection"`
	PrivateConnection  *DatabaseConnection `json:"private_connection,omitempty"`
	Region             string              `json:"region"`
	Status             string              `json:"status"`
	CreatedAt          time.Time           `json:"created_at"`
	PrivateNetworkUUID string              `json:"private_network_uuid,omitempty"`
}

// DatabasePool represents a database connection pool
type DatabasePool struct {
	User              string              `json:"user"`
	Name              string              `json:"name"`
	Size              int                 `json:"size"`
	Database          string              `json:"db"`
	Mode              string              `json:"mode"`
	Connection        *DatabaseConnection `json:"connection"`
	PrivateConnection *DatabaseConnection `json:"private_connection,omitempty"`
}

// DatabaseCreatePoolRequest is used to create a new database connection pool
type DatabaseCreatePoolRequest struct {
	User     string `json:"user"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
	Database string `json:"db"`
	Mode     string `json:"mode"`
}

// DatabaseCreateUserRequest is used to create a new database user
type DatabaseCreateUserRequest struct {
	Name string `json:"name"`
}

// DatabaseCreateDBRequest is used to create a new engine-specific database within the cluster
type DatabaseCreateDBRequest struct {
	Name string `json:"name"`
}

// DatabaseCreateReplicaRequest is used to create a new read-only replica
type DatabaseCreateReplicaRequest struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	Size   string `json:"size"`
}

type databaseUserRoot struct {
	User *DatabaseUser `json:"user"`
}

type databaseUsersRoot struct {
	Users []DatabaseUser `json:"users"`
}

type databaseDBRoot struct {
	DB *DatabaseDB `json:"db"`
}

type databaseDBsRoot struct {
	DBs []DatabaseDB `json:"dbs"`
}

type databasesRoot struct {
	Databases []Database `json:"databases"`
}

type databaseRoot struct {
	Database *Database `json:"database"`
}

type databaseBackupsRoot struct {
	Backups []DatabaseBackup `json:"backups"`
}

type databasePoolRoot struct {
	Pool *DatabasePool `json:"pool"`
}

type databasePoolsRoot struct {
	Pools []DatabasePool `json:"pools"`
}

type databaseReplicaRoot struct {
	Replica *DatabaseReplica `json:"replica"`
}

type databaseReplicasRoot struct {
	Replicas []DatabaseReplica `json:"replicas"`
}

type evictionPolicyRoot struct {
	EvictionPolicy string `json:"eviction_policy"`
}

// List returns a list of the Databases visible with the caller's API token
func (svc *DatabasesServiceOp) List(ctx context.Context, opts *ListOptions) ([]Database, *Response, error) {
	path := databaseBasePath
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databasesRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Databases, resp, nil
}

// Get retrieves the details of a database cluster
func (svc *DatabasesServiceOp) Get(ctx context.Context, databaseID string) (*Database, *Response, error) {
	path := fmt.Sprintf(databaseSinglePath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Database, resp, nil
}

// Create creates a database cluster
func (svc *DatabasesServiceOp) Create(ctx context.Context, create *DatabaseCreateRequest) (*Database, *Response, error) {
	path := databaseBasePath
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Database, resp, nil
}

// Delete deletes a database cluster. There is no way to recover a cluster once
// it has been destroyed.
func (svc *DatabasesServiceOp) Delete(ctx context.Context, databaseID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", databaseBasePath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// Resize resizes a database cluster by number of nodes or size
func (svc *DatabasesServiceOp) Resize(ctx context.Context, databaseID string, resize *DatabaseResizeRequest) (*Response, error) {
	path := fmt.Sprintf(databaseResizePath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, resize)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// Migrate migrates a database cluster to a new region
func (svc *DatabasesServiceOp) Migrate(ctx context.Context, databaseID string, migrate *DatabaseMigrateRequest) (*Response, error) {
	path := fmt.Sprintf(databaseMigratePath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, migrate)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// UpdateMaintenance updates the maintenance window on a cluster
func (svc *DatabasesServiceOp) UpdateMaintenance(ctx context.Context, databaseID string, maintenance *DatabaseUpdateMaintenanceRequest) (*Response, error) {
	path := fmt.Sprintf(databaseMaintenancePath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, maintenance)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// ListBackups returns a list of the current backups of a database
func (svc *DatabasesServiceOp) ListBackups(ctx context.Context, databaseID string, opts *ListOptions) ([]DatabaseBackup, *Response, error) {
	path := fmt.Sprintf(databaseBackupsPath, databaseID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseBackupsRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Backups, resp, nil
}

// GetUser returns the database user identified by userID
func (svc *DatabasesServiceOp) GetUser(ctx context.Context, databaseID, userID string) (*DatabaseUser, *Response, error) {
	path := fmt.Sprintf(databaseUserPath, databaseID, userID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseUserRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.User, resp, nil
}

// ListUsers returns all database users for the database
func (svc *DatabasesServiceOp) ListUsers(ctx context.Context, databaseID string, opts *ListOptions) ([]DatabaseUser, *Response, error) {
	path := fmt.Sprintf(databaseUsersPath, databaseID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseUsersRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Users, resp, nil
}

// CreateUser will create a new database user
func (svc *DatabasesServiceOp) CreateUser(ctx context.Context, databaseID string, createUser *DatabaseCreateUserRequest) (*DatabaseUser, *Response, error) {
	path := fmt.Sprintf(databaseUsersPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, createUser)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseUserRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.User, resp, nil
}

// DeleteUser will delete an existing database user
func (svc *DatabasesServiceOp) DeleteUser(ctx context.Context, databaseID, userID string) (*Response, error) {
	path := fmt.Sprintf(databaseUserPath, databaseID, userID)
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// ListDBs returns all databases for a given database cluster
func (svc *DatabasesServiceOp) ListDBs(ctx context.Context, databaseID string, opts *ListOptions) ([]DatabaseDB, *Response, error) {
	path := fmt.Sprintf(databaseDBsPath, databaseID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseDBsRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.DBs, resp, nil
}

// GetDB returns a single database by name
func (svc *DatabasesServiceOp) GetDB(ctx context.Context, databaseID, name string) (*DatabaseDB, *Response, error) {
	path := fmt.Sprintf(databaseDBPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseDBRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.DB, resp, nil
}

// CreateDB will create a new database
func (svc *DatabasesServiceOp) CreateDB(ctx context.Context, databaseID string, createDB *DatabaseCreateDBRequest) (*DatabaseDB, *Response, error) {
	path := fmt.Sprintf(databaseDBsPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, createDB)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseDBRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.DB, resp, nil
}

// DeleteDB will delete an existing database
func (svc *DatabasesServiceOp) DeleteDB(ctx context.Context, databaseID, name string) (*Response, error) {
	path := fmt.Sprintf(databaseDBPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// ListPools returns all connection pools for a given database cluster
func (svc *DatabasesServiceOp) ListPools(ctx context.Context, databaseID string, opts *ListOptions) ([]DatabasePool, *Response, error) {
	path := fmt.Sprintf(databasePoolsPath, databaseID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databasePoolsRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Pools, resp, nil
}

// GetPool returns a single database connection pool by name
func (svc *DatabasesServiceOp) GetPool(ctx context.Context, databaseID, name string) (*DatabasePool, *Response, error) {
	path := fmt.Sprintf(databasePoolPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databasePoolRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Pool, resp, nil
}

// CreatePool will create a new database connection pool
func (svc *DatabasesServiceOp) CreatePool(ctx context.Context, databaseID string, createPool *DatabaseCreatePoolRequest) (*DatabasePool, *Response, error) {
	path := fmt.Sprintf(databasePoolsPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, createPool)
	if err != nil {
		return nil, nil, err
	}
	root := new(databasePoolRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Pool, resp, nil
}

// DeletePool will delete an existing database connection pool
func (svc *DatabasesServiceOp) DeletePool(ctx context.Context, databaseID, name string) (*Response, error) {
	path := fmt.Sprintf(databasePoolPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetReplica returns a single database replica
func (svc *DatabasesServiceOp) GetReplica(ctx context.Context, databaseID, name string) (*DatabaseReplica, *Response, error) {
	path := fmt.Sprintf(databaseReplicaPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseReplicaRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Replica, resp, nil
}

// ListReplicas returns all read-only replicas for a given database cluster
func (svc *DatabasesServiceOp) ListReplicas(ctx context.Context, databaseID string, opts *ListOptions) ([]DatabaseReplica, *Response, error) {
	path := fmt.Sprintf(databaseReplicasPath, databaseID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseReplicasRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Replicas, resp, nil
}

// CreateReplica will create a new database connection pool
func (svc *DatabasesServiceOp) CreateReplica(ctx context.Context, databaseID string, createReplica *DatabaseCreateReplicaRequest) (*DatabaseReplica, *Response, error) {
	path := fmt.Sprintf(databaseReplicasPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodPost, path, createReplica)
	if err != nil {
		return nil, nil, err
	}
	root := new(databaseReplicaRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Replica, resp, nil
}

// DeleteReplica will delete an existing database replica
func (svc *DatabasesServiceOp) DeleteReplica(ctx context.Context, databaseID, name string) (*Response, error) {
	path := fmt.Sprintf(databaseReplicaPath, databaseID, name)
	req, err := svc.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetEvictionPolicy loads the eviction policy for a given Redis cluster.
func (svc *DatabasesServiceOp) GetEvictionPolicy(ctx context.Context, databaseID string) (string, *Response, error) {
	path := fmt.Sprintf(evictionPolicyPath, databaseID)
	req, err := svc.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return "", nil, err
	}
	root := new(evictionPolicyRoot)
	resp, err := svc.client.Do(ctx, req, root)
	if err != nil {
		return "", resp, err
	}
	return root.EvictionPolicy, resp, nil
}

// SetEvictionPolicy updates the eviction policy for a given Redis cluster.
func (svc *DatabasesServiceOp) SetEvictionPolicy(ctx context.Context, databaseID, policy string) (*Response, error) {
	path := fmt.Sprintf(evictionPolicyPath, databaseID)
	root := &evictionPolicyRoot{EvictionPolicy: policy}
	req, err := svc.client.NewRequest(ctx, http.MethodPut, path, root)
	if err != nil {
		return nil, err
	}
	resp, err := svc.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
