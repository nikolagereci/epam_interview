package company_test

//
//import (
//	"context"
//	"github.com/gocql/gocql"
//	"github.com/google/uuid"
//	"github.com/ngereci/xm_interview/company"
//	"github.com/stretchr/testify/assert"
//	"testing"
//	"time"
//)
//
//func TestEmbeddedCassandra(t *testing.T) {
//	// Start an embedded Cassandra cluster with 1 node
//	cluster := gocqltest.NewCluster(1)
//
//	// Start the cluster
//	if err := cluster.Start(); err != nil {
//		t.Fatalf("failed to start embedded Cassandra: %v", err)
//	}
//	defer cluster.Stop()
//
//	// Create a session to the embedded Cassandra cluster
//	clusterConfig := gocql.ClusterConfig{
//		Hosts: []string{cluster.IP()},
//		Port:  cluster.Port(),
//	}
//	session, err := clusterConfig.CreateSession()
//	if err != nil {
//		t.Fatalf("failed to create session to embedded Cassandra: %v", err)
//	}
//	defer session.Close()
//
//	// Run your tests using the session to the embedded Cassandra cluster
//	// ...
//
//	// Sleep for a bit to give the Cassandra cluster time to shutdown
//	time.Sleep(1 * time.Second)
//}
//
//func TestCompanyRepository_Create(t *testing.T) {
//	session := createTestSession()
//	defer session.Close()
//
//	// create a test company
//	newCompany := &company.Company{
//		ID:          uuid.New(),
//		Name:        "Test Company",
//		Description: "A test company",
//		Employees:   100,
//		Registered:  true,
//		Type:        company.Cooperative,
//	}
//
//	repo := company.NewRepository(session)
//	err := repo.Create(context.Background(), newCompany)
//
//	// assert that the error is nil
//	assert.NoError(t, err)
//}
//
//func TestCompanyRepository_GetByID(t *testing.T) {
//	session := createTestSession()
//	defer session.Close()
//
//	// create a test company
//	newCompany := &company.Company{
//		ID:          uuid.New(),
//		Name:        "Test Company",
//		Description: "A test company",
//		Employees:   100,
//		Registered:  true,
//		Type:        company.Corporation,
//	}
//
//	repo := company.NewRepository(session)
//	err := repo.Create(context.Background(), newCompany)
//
//	// assert that the error is nil
//	assert.NoError(t, err)
//
//	// get the company by ID
//	result, err := repo.GetByID(context.Background(), company.ID)
//
//	// assert that the error is nil
//	assert.NoError(t, err)
//
//	// assert that the result is not nil
//	assert.NotNil(t, result)
//
//	// assert that the result matches the original company
//	assert.Equal(t, company.ID, result.ID)
//	assert.Equal(t, company.Name, result.Name)
//	assert.Equal(t, company.Description, result.Description)
//	assert.Equal(t, company.Employees, result.Employees)
//	assert.Equal(t, company.Registered, result.Registered)
//	assert.Equal(t, company.Type, result.Type)
//}
