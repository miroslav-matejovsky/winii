package fileinfo

import (
	"crypto/x509"
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Certificates struct {
	Certificates []*x509.Certificate
}

func (c *Certificates) SignedBy(verifier string) bool {
	if len(c.Certificates) == 0 {
		return false
	}
	for _, cert := range c.Certificates {
		if strings.EqualFold(cert.Subject.CommonName, verifier) {
			return true
		}
	}
	return false
}

func (c *Certificates) ValidAtSignedBy(verifier string, at time.Time) bool {
	if len(c.Certificates) == 0 {
		return false
	}
	for _, cert := range c.Certificates {
		valid := cert.NotBefore.Before(at) && cert.NotAfter.After(at)
		if strings.EqualFold(cert.Subject.CommonName, verifier) && valid {
			return true
		}
	}
	return false
}

func getCertificates(filepath string) (*Certificates, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	peFile, err := pe.NewFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PE file: %v", err)
	}
	defer func() {
		_ = peFile.Close()
	}()

	// Extract certificates from the PE file
	certs, err := extractCertificates(file, peFile)
	if err != nil {
		return nil, fmt.Errorf("failed to extract certificates: %v", err)
	}

	return &Certificates{Certificates: certs}, nil
}

func extractCertificates(file *os.File, peFile *pe.File) ([]*x509.Certificate, error) {
	var certDir pe.DataDirectory

	// Get the certificate table from the data directory
	switch oh := peFile.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		if len(oh.DataDirectory) <= 4 {
			return nil, fmt.Errorf("no certificate table found")
		}
		certDir = oh.DataDirectory[4] // IMAGE_DIRECTORY_ENTRY_SECURITY
	case *pe.OptionalHeader64:
		if len(oh.DataDirectory) <= 4 {
			return nil, fmt.Errorf("no certificate table found")
		}
		certDir = oh.DataDirectory[4] // IMAGE_DIRECTORY_ENTRY_SECURITY
	default:
		return nil, fmt.Errorf("unsupported PE format")
	}

	if certDir.Size == 0 {
		return nil, nil
	}

	// Read the certificate table
	_, err := file.Seek(int64(certDir.VirtualAddress), 0)
	if err != nil {
		return nil, fmt.Errorf("failed to seek to certificate table: %v", err)
	}

	certData := make([]byte, certDir.Size)
	_, err = io.ReadFull(file, certData)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate data: %v", err)
	}

	return parseCertificateTable(certData)
}

func parseCertificateTable(data []byte) ([]*x509.Certificate, error) {
	var certificates []*x509.Certificate
	offset := 0

	for offset < len(data) {
		if offset+8 > len(data) {
			break
		}

		// Read the WIN_CERTIFICATE structure
		length := binary.LittleEndian.Uint32(data[offset:])
		// revision := binary.LittleEndian.Uint16(data[offset+4:])
		certType := binary.LittleEndian.Uint16(data[offset+6:])

		if length < 8 || offset+int(length) > len(data) {
			break
		}

		// slog.Debug("Found certificate entry",
		// 	"Length", length,
		// 	"Revision", fmt.Sprintf("0x%x", revision),
		// 	"Type", fmt.Sprintf("0x%x", certType))

		// WIN_CERT_TYPE_PKCS_SIGNED_DATA = 0x0002
		if certType == 0x0002 {
			certBytes := data[offset+8 : offset+int(length)]
			certs := extractCertsFromPKCS7(certBytes)
			certificates = append(certificates, certs...)
		}

		// Move to next certificate (aligned to 8-byte boundary)
		offset += int(length)
		offset = (offset + 7) &^ 7
	}

	return certificates, nil
}

func extractCertsFromPKCS7(data []byte) []*x509.Certificate {
	var certificates []*x509.Certificate

	// Try to parse as PKCS#7 DER format
	// This is a simplified approach - in practice, you'd want to properly parse the ASN.1 structure

	// Look for certificate patterns in the PKCS#7 data
	// X.509 certificates in DER format start with 0x30 0x82 (SEQUENCE tag)
	for i := 0; i < len(data)-4; i++ {
		if data[i] == 0x30 && data[i+1] == 0x82 {
			// Try to parse certificate starting at this position
			remaining := data[i:]

			// Parse the length field
			if len(remaining) < 4 {
				continue
			}

			certLen := int(remaining[2])<<8 + int(remaining[3]) + 4
			if certLen > len(remaining) {
				continue
			}

			certData := remaining[:certLen]
			cert, err := x509.ParseCertificate(certData)
			if err == nil {
				certificates = append(certificates, cert)
				i += certLen - 1 // Skip past this certificate
			}
		}
	}

	// If the simple approach doesn't work, try parsing the entire block
	if len(certificates) == 0 {
		certs, err := x509.ParseCertificates(data)
		if err == nil {
			certificates = append(certificates, certs...)
		}
	}

	return certificates
}
