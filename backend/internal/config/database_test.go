package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractDatabaseNameFromDSN(t *testing.T) {
	dbName, err := extractDatabaseName("host=localhost user=taskuser password=secret dbname=taskdb_test port=5432 sslmode=disable")
	require.NoError(t, err)
	require.Equal(t, "taskdb_test", dbName)
}

func TestExtractDatabaseNameFromDSN_Empty(t *testing.T) {
	_, err := extractDatabaseName("")
	require.Error(t, err)
}
