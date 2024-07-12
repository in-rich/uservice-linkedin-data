package dao_test

import (
	"context"
	"github.com/in-rich/uservice-linkedin-data/pkg/dao"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var getProfilePictureFixtures = map[string][]byte{
	"public-identifier-1": []byte(RawSus),
}

func TestGetProfilePicture(t *testing.T) {
	bucket := NewStorage("user-profile-pictures-test", getProfilePictureFixtures)
	defer ClearStorage(bucket)
	repository := dao.NewGetProfilePictureRepository(bucket)

	testData := []struct {
		name string

		publicIdentifier string

		expect        string
		expectContent string
		expectErr     error
	}{
		{
			name:             "GetProfilePicture",
			publicIdentifier: "public-identifier-1",
			expect:           StorageURI("public-identifier-1"),
			expectContent:    RawSus,
		},
		{
			name:             "ProfilePictureNotFound",
			publicIdentifier: "public-identifier-2",
			expectErr:        dao.ErrProfilePictureNotFound,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			url, err := repository.GetProfilePicture(context.TODO(), tt.publicIdentifier)

			require.ErrorIs(t, err, tt.expectErr)
			require.True(t, strings.HasPrefix(url, tt.expect))

			if tt.expectErr == nil {
				content := DownloadBase64(url)
				require.Equal(t, tt.expectContent, content)
			}
		})
	}
}
