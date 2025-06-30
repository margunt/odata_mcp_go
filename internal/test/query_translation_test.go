package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/odata-mcp/go/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestQueryParameterTranslation tests v2 to v4 query parameter translation
func TestQueryParameterTranslation(t *testing.T) {
	tests := []struct {
		name            string
		isV4            bool
		inputOptions    map[string]string
		expectedQueries map[string]string // Expected query parameters in URL
		notExpected     []string          // Parameters that should NOT be in URL
	}{
		{
			name: "V2 service with inlinecount",
			isV4: false,
			inputOptions: map[string]string{
				"$filter": "Name eq 'Test'",
				"$top":    "10",
			},
			expectedQueries: map[string]string{
				"$format":      "json",
				"$inlinecount": "allpages",
				"$filter":      "Name eq 'Test'",
				"$top":         "10",
			},
			notExpected: []string{"$count"},
		},
		{
			name: "V4 service without explicit count",
			isV4: true,
			inputOptions: map[string]string{
				"$filter": "contains(Name, 'Test')",
				"$top":    "10",
			},
			expectedQueries: map[string]string{
				"$filter": "contains(Name, 'Test')",
				"$top":    "10",
			},
			notExpected: []string{"$format", "$inlinecount", "$count"},
		},
		{
			name: "V4 service with explicit inlinecount translation",
			isV4: true,
			inputOptions: map[string]string{
				"$filter":      "contains(Name, 'Test')",
				"$top":         "10",
				"$inlinecount": "allpages",
			},
			expectedQueries: map[string]string{
				"$filter": "contains(Name, 'Test')",
				"$top":    "10",
				"$count":  "true",
			},
			notExpected: []string{"$format", "$inlinecount"},
		},
		{
			name: "V4 service with inlinecount=none translation",
			isV4: true,
			inputOptions: map[string]string{
				"$filter":      "ProductID gt 5",
				"$inlinecount": "none",
			},
			expectedQueries: map[string]string{
				"$filter": "ProductID gt 5",
				"$count":  "false",
			},
			notExpected: []string{"$format", "$inlinecount"},
		},
		{
			name: "V2 service with explicit count (should keep inlinecount)",
			isV4: false,
			inputOptions: map[string]string{
				"$filter": "Price gt 100",
				"$count":  "true", // This should be kept as-is for v2
			},
			expectedQueries: map[string]string{
				"$format":      "json",
				"$inlinecount": "allpages",
				"$filter":      "Price gt 100",
				"$count":       "true",
			},
			notExpected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedURL *url.URL
			
			// Create test server that captures the request URL
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedURL = r.URL
				
				// Return appropriate metadata based on version
				if r.URL.Path == "/$metadata" {
					w.Header().Set("Content-Type", "application/xml")
					if tt.isV4 {
						w.Write([]byte(`<?xml version="1.0" encoding="utf-8"?>
<edmx:Edmx Version="4.0" xmlns:edmx="http://docs.oasis-open.org/odata/ns/edmx">
  <edmx:DataServices>
    <Schema Namespace="Test" xmlns="http://docs.oasis-open.org/odata/ns/edm">
      <EntityContainer Name="TestContainer">
        <EntitySet Name="TestEntities" EntityType="Test.TestEntity" />
      </EntityContainer>
    </Schema>
  </edmx:DataServices>
</edmx:Edmx>`))
					} else {
						w.Write([]byte(`<?xml version="1.0" encoding="utf-8"?>
<edmx:Edmx Version="1.0" xmlns:edmx="http://schemas.microsoft.com/ado/2007/06/edmx">
  <edmx:DataServices m:DataServiceVersion="2.0" xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata">
    <Schema Namespace="Test" xmlns="http://schemas.microsoft.com/ado/2009/11/edm">
      <EntityContainer Name="TestContainer" m:IsDefaultEntityContainer="true">
        <EntitySet Name="TestEntities" EntityType="Test.TestEntity" />
      </EntityContainer>
    </Schema>
  </edmx:DataServices>
</edmx:Edmx>`))
					}
					return
				}
				
				// Return empty collection for entity queries
				w.Header().Set("Content-Type", "application/json")
				if tt.isV4 {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"@odata.context": "http://" + r.Host + "/$metadata#TestEntities",
						"value":          []interface{}{},
					})
				} else {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"d": map[string]interface{}{
							"results": []interface{}{},
						},
					})
				}
			}))
			defer server.Close()
			
			// Create client
			odataClient := client.NewODataClient(server.URL, false)
			ctx := context.Background()
			
			// Fetch metadata to set version
			_, err := odataClient.GetMetadata(ctx)
			require.NoError(t, err)
			
			// Make request with options
			_, err = odataClient.GetEntitySet(ctx, "TestEntities", tt.inputOptions)
			require.NoError(t, err)
			
			// Verify query parameters
			require.NotNil(t, capturedURL)
			queryParams := capturedURL.Query()
			
			// Check expected parameters
			for param, expectedValue := range tt.expectedQueries {
				actualValue := queryParams.Get(param)
				assert.Equal(t, expectedValue, actualValue, 
					"Expected %s=%s but got %s", param, expectedValue, actualValue)
			}
			
			// Check parameters that should NOT be present
			for _, param := range tt.notExpected {
				assert.Empty(t, queryParams.Get(param), 
					"Parameter %s should not be present in URL", param)
			}
		})
	}
}

// TestCountToolTranslation tests that count tools work correctly with both v2 and v4
func TestCountToolTranslation(t *testing.T) {
	tests := []struct {
		name         string
		isV4         bool
		responseBody string
		expectedURL  string
	}{
		{
			name: "V2 count request",
			isV4: false,
			responseBody: `{"d": {"results": [], "__count": "42"}}`,
			expectedURL: "TestEntities?%24format=json&%24inlinecount=allpages&%24top=0",
		},
		{
			name: "V4 count request",
			isV4: true,
			responseBody: `{"@odata.context": "http://example.com/$metadata#TestEntities", "@odata.count": 42, "value": []}`,
			expectedURL: "TestEntities?%24count=true&%24top=0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedPath string
			
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/$metadata" {
					w.Header().Set("Content-Type", "application/xml")
					if tt.isV4 {
						w.Write([]byte(`<?xml version="1.0" encoding="utf-8"?>
<edmx:Edmx Version="4.0" xmlns:edmx="http://docs.oasis-open.org/odata/ns/edmx">
  <edmx:DataServices>
    <Schema Namespace="Test" xmlns="http://docs.oasis-open.org/odata/ns/edm">
      <EntityContainer Name="TestContainer">
        <EntitySet Name="TestEntities" EntityType="Test.TestEntity" />
      </EntityContainer>
    </Schema>
  </edmx:DataServices>
</edmx:Edmx>`))
					} else {
						w.Write([]byte(`<?xml version="1.0" encoding="utf-8"?>
<edmx:Edmx Version="1.0" xmlns:edmx="http://schemas.microsoft.com/ado/2007/06/edmx">
  <edmx:DataServices m:DataServiceVersion="2.0" xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata">
    <Schema Namespace="Test" xmlns="http://schemas.microsoft.com/ado/2009/11/edm">
      <EntityContainer Name="TestContainer" m:IsDefaultEntityContainer="true">
        <EntitySet Name="TestEntities" EntityType="Test.TestEntity" />
      </EntityContainer>
    </Schema>
  </edmx:DataServices>
</edmx:Edmx>`))
					}
					return
				}
				
				capturedPath = r.URL.Path + "?" + r.URL.RawQuery
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()
			
			// Create client
			odataClient := client.NewODataClient(server.URL, false)
			ctx := context.Background()
			
			// Fetch metadata
			_, err := odataClient.GetMetadata(ctx)
			require.NoError(t, err)
			
			// Simulate count tool request
			options := map[string]string{
				"$inlinecount": "allpages",
				"$top":         "0",
			}
			
			resp, err := odataClient.GetEntitySet(ctx, "TestEntities", options)
			require.NoError(t, err)
			
			// Verify URL
			assert.Contains(t, capturedPath, tt.expectedURL)
			
			// Verify count was parsed correctly
			assert.NotNil(t, resp.Count)
			assert.Equal(t, int64(42), *resp.Count)
		})
	}
}