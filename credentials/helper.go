package credentials

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/zhangel/go-frame.git/certificate"
	"github.com/zhangel/go-frame.git/config"
	"github.com/zhangel/go-frame.git/config/watcher"
	"github.com/zhangel/go-frame.git/log"
	"github.com/google/uuid"
)

type onProduceCertificateFn func(certificate.Request, *tls.Certificate, error)

func (s onProduceCertificateFn) Consume(certInfo certificate.Request, certificate *tls.Certificate, err error) {
	s(certInfo, certificate, err)
}

type IntermediateCAProvider interface {
	GetIntermediateCA(ctx context.Context) (cert, key []byte, err error)
	WatchIntermediateCA(func(cert, key []byte)) (canceler func(), err error)
}

type certRequestFromConfig struct {
	id             string
	commonNameKeys []string

	mu          sync.RWMutex
	commonNames []string
}

func NewCertRequestFromConfig(commonNameKeys ...string) certificate.Request {
	var id string
	if uuid, err := uuid.NewUUID(); err != nil {
		id = "default"
	} else {
		id = uuid.String()
	}

	req := &certRequestFromConfig{
		id:             id,
		commonNameKeys: commonNameKeys,
	}
	req.update()

	return req
}

func (s *certRequestFromConfig) Id() string {
	return s.id
}

func (s *certRequestFromConfig) update() {
	commonNameSet := make(map[string]struct{})
	insertCommonNames := func(commonNames []string) {
		for _, commonName := range commonNames {
			commonName := strings.TrimSpace(commonName)
			if commonName == "" {
				continue
			}
			commonNameSet[commonName] = struct{}{}
		}
	}

	for _, key := range s.commonNameKeys {
		insertCommonNames(config.StringList(key))
	}
	insertCommonNames(GenerateSAN(nil))

	s.mu.Lock()
	defer s.mu.Unlock()

	commonNames := make([]string, 0, len(commonNameSet))
	for commonName := range commonNameSet {
		commonNames = append(commonNames, commonName)
	}
	sort.Slice(commonNames, func(i, j int) bool {
		return commonNames[i] > commonNames[j]
	})
	s.commonNames = commonNames
}

func (s *certRequestFromConfig) CommonName() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.commonNames) > 0 {
		return s.commonNames[0]
	} else {
		return ""
	}
}

func (s *certRequestFromConfig) IpSans() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.commonNames
}

func (s *certRequestFromConfig) DomainSans() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.commonNames
}

func (s *certRequestFromConfig) CertChain() bool {
	return true
}

func (s *certRequestFromConfig) PubkeyBits() int {
	return 2048
}

func (s *certRequestFromConfig) RegisterNotify(notify func()) (func(), error) {
	var cancelers []func()
	for _, key := range s.commonNameKeys {
		cancelers = append(cancelers, config.Watch(watcher.NewHelper(key, func(v string, deleted bool) {
			s.update()
			if notify != nil {
				notify()
			}
		})))
	}

	return func() {
		for _, canceler := range cancelers {
			if canceler != nil {
				canceler()
			}
		}
	}, nil
}

func (s *certRequestFromConfig) String() string {
	return fmt.Sprintf("CN = %s, DNSNames = %+v, IPAddresses = %+v, CertChain = %v", s.CommonName(), s.DomainSans(), s.IpSans(), s.CertChain())
}

func NewCertificateFactory(intermediateCAProvider IntermediateCAProvider, certRequest certificate.Request, updateNotify func()) (factory *certificate.Factory, watchCanceler func(), err error) {
	if intermediateCAProvider == nil {
		return nil, nil, fmt.Errorf("NewCertificateFactory failed, err = intermediate CA provider is nil")
	}

	var mu sync.RWMutex
	var currentCert []byte
	var currentKey []byte
	var newCert []byte
	var newKey []byte

	issuer, err := certificate.NewIssuer(
		certificate.IssuerWithCAProvider(func() (caCert, caKey, passphrase []byte, err error) {
			mu.Lock()
			defer mu.Unlock()

			currentCert, currentKey, err = intermediateCAProvider.GetIntermediateCA(context.Background())
			newCert = currentCert
			newKey = currentKey
			return currentCert, currentKey, nil, err
		}),
	)

	if err != nil {
		log.Errorf("NewCertificateFactory, create certificate issuer failed, err = %v", err)
		return nil, nil, err
	}

	factory, err = certificate.NewCertificateFactory(issuer, certificate.NewDefaultStore())
	if err != nil {
		log.Errorf("NewCertificateFactory, create certificate factory failed, err = %v", err)
		return nil, nil, err
	}
	factory.Produce(certRequest)

	var cancelers []func()
	if canceler, err := intermediateCAProvider.WatchIntermediateCA(func(cert, key []byte) {
		mu.Lock()

		if len(cert) != 0 {
			newCert = cert
		}

		if len(key) != 0 {
			newKey = key
		}

		if !bytes.Equal(newCert, currentCert) && !bytes.Equal(newKey, currentKey) {
			currentCert = newCert
			currentKey = newKey
			log.Info("Renewal server certificate with new CA certificate...")

			mu.Unlock()
			factory.Produce(certRequest)
		} else {
			mu.Unlock()
		}
	}); err != nil {
		log.Errorf("NewCertificateFactory, watch intermediateCA failed, err = %v", err)
		return nil, nil, err
	} else {
		cancelers = append(cancelers, canceler)
	}

	if reqChangeNotify, ok := certRequest.(certificate.RequestChangeNotify); ok {
		if canceler, err := reqChangeNotify.RegisterNotify(func() {
			log.Infof("Renewal server certificate with new certificate request [%v]...", certRequest)
			factory.Produce(certRequest)
		}); err == nil {
			cancelers = append(cancelers, canceler)
		} else {
			log.Errorf("CertificateFactory, register common name change notify failed, err = %+v", err)
			return nil, nil, err
		}
	}

	factory.RegisterConsumer(onProduceCertificateFn(func(r certificate.Request, c *tls.Certificate, err error) {
		if updateNotify != nil {
			updateNotify()
		}
	}))

	return factory, func() {
		for _, canceler := range cancelers {
			if canceler != nil {
				canceler()
			}
		}
	}, nil
}
