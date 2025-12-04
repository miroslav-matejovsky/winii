package fileinfo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNonExistentFile(t *testing.T) {
	_, err := NewWinFileInfo("C:\\nonexistent.txt")
	require.ErrorContains(t, err, "file does not exist")
}

const signer = "some signer"

func TestSignedExe(t *testing.T) {
	t.Skip("no test data")
	wfi, err := NewWinFileInfo("testdata\\signed.exe")
	require.NoError(t, err)
	certificates, err := wfi.GetCertificates()
	require.NoError(t, err)
	require.NotNil(t, certificates)

	require.NotEmpty(t, certificates.Certificates, "Signed executable should have at least one certificate")

	valid := certificates.SignedBy(signer)
	require.True(t, valid, "Signed executable should be signed")

	validAt := certificates.ValidAtSignedBy(signer, time.Now())
	require.True(t, validAt, "Signed executable should be valid at the current time")

	versions, err := wfi.GetVersions()
	require.NoError(t, err)
	require.NotNil(t, versions)

	require.Equal(t, "5.3.0.0", versions.FileVersion.String())
}

func TestInstalledExe(t *testing.T) {
	t.Skip("no test data")

	wfi, err := NewWinFileInfo("testdata\\installed.exe")
	require.NoError(t, err)
	certificates, err := wfi.GetCertificates()
	require.NoError(t, err)
	require.NotNil(t, certificates)

	require.NotEmpty(t, certificates.Certificates, "Signed executable should have at least one certificate")

	valid := certificates.SignedBy(signer)
	require.True(t, valid, "Signed executable should be signed")

	validAt := certificates.ValidAtSignedBy(signer, time.Now())
	require.True(t, validAt, "Signed executable should be valid at the current time")

	versions, err := wfi.GetVersions()
	require.NoError(t, err)
	require.NotNil(t, versions)

	require.Equal(t, "5.3.0.0", versions.FileVersion.String())
}

func TestUnsignedExe(t *testing.T) {
	t.Skip("no test data")

	wfi, err := NewWinFileInfo("testdata\\unsigned.exe")
	require.NoError(t, err)
	certificates, err := wfi.GetCertificates()
	require.NoError(t, err)
	require.NotNil(t, certificates)

	require.Empty(t, certificates.Certificates, "Unsigned executable should not have any certificates")

	versions, err := wfi.GetVersions()
	require.NoError(t, err)
	require.NotNil(t, versions)

	require.Equal(t, "5.3.0.0", versions.FileVersion.String())
}
