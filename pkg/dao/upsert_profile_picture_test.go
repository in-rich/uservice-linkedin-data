package dao_test

import (
	"context"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var upsertProfilePictureFixtures = map[string][]byte{
	"public-identifier-1": []byte("foo"),
}

func TestUpsertProfilePicture(t *testing.T) {
	testData := []struct {
		name string

		publicIdentifier string
		base64           string

		expect        string
		expectContent string
		expectErr     error
	}{
		{
			name:             "UpdateProfilePicture",
			publicIdentifier: "public-identifier-1",
			base64:           Base64Sus,
			expect:           StorageURI("public-identifier-1"),
			expectContent:    RawSus,
		},
		{
			name:             "CreateProfilePicture",
			publicIdentifier: "public-identifier-2",
			base64:           Base64Sus,
			expect:           StorageURI("public-identifier-2"),
			expectContent:    RawSus,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			bucket := NewStorage("user-profile-pictures-test", upsertProfilePictureFixtures)
			defer ClearStorage(bucket)

			repository := dao.NewUpsertProfilePictureRepository(bucket)

			url, err := repository.UpsertProfilePicture(context.TODO(), tt.publicIdentifier, tt.base64)

			require.NoError(t, err)
			require.True(t, strings.HasPrefix(url, tt.expect))

			if tt.expectErr == nil {
				content := DownloadBase64(url)
				require.Equal(t, tt.expectContent, content)
			}
		})
	}
}
