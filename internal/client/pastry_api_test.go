/*
 * Copyright The Microcks Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package client_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/microcks/microcks-testcontainers-go-demo/internal/client"
	microcks "microcks.io/testcontainers-go"
)

func setup(ctx context.Context, t *testing.T) *microcks.MicrocksContainer {
	microcksContainer, err := microcks.Run(ctx,
		"quay.io/microcks/microcks-uber:1.10.0-native",
		microcks.WithMainArtifact("../../testdata/apipastries-openapi.yaml"),
		microcks.WithSecondaryArtifact("../../testdata/apipastries-postman-collection.json"),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		if err := microcksContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})
	return microcksContainer
}

func TestGetPastry(t *testing.T) {
	ctx := context.Background()
	microcksContainer := setup(ctx, t)

	baseApiUrl, err := microcksContainer.RestMockEndpoint(ctx, "API Pastries", "0.0.1")
	require.NoError(t, err)
	pastryAPIClient := client.NewPastryAPIClient(baseApiUrl)

	pastry, err := pastryAPIClient.GetPastry("Millefeuille")
	require.NoError(t, err)
	require.Equal(t, "Millefeuille", pastry.Name)
	require.Equal(t, "available", pastry.Status)

	pastry, err = pastryAPIClient.GetPastry("Eclair Cafe")
	require.NoError(t, err)
	require.Equal(t, "Eclair Cafe", pastry.Name)
	require.Equal(t, "available", pastry.Status)

	pastry, err = pastryAPIClient.GetPastry("Eclair Chocolat")
	require.NoError(t, err)
	require.Equal(t, "Eclair Chocolat", pastry.Name)
	require.Equal(t, "unknown", pastry.Status)
}

func TestListPastries(t *testing.T) {
	ctx := context.Background()
	microcksContainer := setup(ctx, t)

	baseApiUrl, err := microcksContainer.RestMockEndpoint(ctx, "API Pastries", "0.0.1")
	require.NoError(t, err)
	pastryAPIClient := client.NewPastryAPIClient(baseApiUrl)

	pastries, err := pastryAPIClient.ListPastries("S")
	require.NoError(t, err)
	require.Equal(t, 1, len(*pastries))

	pastries, err = pastryAPIClient.ListPastries("M")
	require.NoError(t, err)
	require.Equal(t, 2, len(*pastries))

	pastries, err = pastryAPIClient.ListPastries("L")
	require.NoError(t, err)
	require.Equal(t, 2, len(*pastries))
}
