package certificate

const selfSignedCACert = `-----BEGIN CERTIFICATE-----
MIIFJjCCAw6gAwIBAgIBATANBgkqhkiG9w0BAQsFADAyMQwwCgYDVQQKEwNRQVgx
DDAKBgNVBAsTA0FUUzEUMBIGA1UEAxMLR29GcmFtZXdvcmswIBcNMjAwMTMxMTAz
NDQ3WhgPMjA3MDAxMzExMDM0NDRaMDIxDDAKBgNVBAoTA1FBWDEMMAoGA1UECxMD
QVRTMRQwEgYDVQQDEwtHb0ZyYW1ld29yazCCAiIwDQYJKoZIhvcNAQEBBQADggIP
ADCCAgoCggIBAOJiaiJ01afeKHiiDETM12XXn9HxM50LrQAAYJU6U4271qEpIOUx
mUHBIgu8N3LccKp2Tovzgz8vPU6L9d058EFdvH9SCg+Zj3GlU7xrzba7GBUKe9CT
NICEaYiZsNKzh4ZiwABXwOg1EWSNMxc8nXR7qeEmMGIhTEW6vIIJlXnun2lebfgY
LvpP2GsJCo+nx4Bgr1eZXEndvHkru5tX1rh6K2cSu6xB9cOWG15BMNyO+ZusXihK
fGmxnXDq7PkdycumJExCRZS5cSkEcijQ03P0HUjoSveORbkbpUX81dLB9KAjzM/i
sIuKiRiZXLPNyEhsolz4DDJivXcGD0B6EX2Ef+vNcObmR4rvG+Zsb2pVME2Vvidr
QMx+3xpKA6JZFriKIpGANWrm+ilAARjY2rsGAxHiiiamWn2sP+tIGKtMzxXvlHxS
l/UjvhHqxuuxQbWy0G+9cA7Rtj2Hp+LfIqNBPKuA1I+aUvUqSrXYMfo4eapKvO2/
Ycb+IYC22qE5Y8L7Dpk8g5HTaI1g2LrwUHH462aUgT6zmqwyLmDJBUahqbBstzVH
C5c4TW8kfgCrLX0WVZDq5dZupmteQS2BiDUctbmN6DaSHyl1EVlZ1CeRu17FycAh
YlKbJmy7+EkidVv80OoeqXL7UGwRrfBLPtVEF00S3jh2M5UCX84a188lAgMBAAGj
RTBDMA4GA1UdDwEB/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/AgEBMB0GA1UdDgQW
BBSVcX8zPNEN/wed9gA9RyBFsHVDdjANBgkqhkiG9w0BAQsFAAOCAgEAJWQiz2oH
FXNU2T8g/x1UaNlOQcH7H+QYyZWaNpZW5P66OMO6tE4bfXm9mPzSOi2nGNkzpLd4
u+ok3Dz6liwHme2eZg8PYmMOiBauSrY71y4A1wK2+MIJyq6+xtqAnnf5JWLktUc5
REpFJqltHMgQYsld7tOnTpyatW9XbaCT0okzdxykfg5QAmBFRYnzzLX5/U54HNlP
VML43zQig0w+BmHVt/kmg9hCVUt/DpQHWyeQ+Gk5wi/Hmvkpz812kf0nrVlNTa/E
+XjrD+Zib/q7ANufhHchoWzcUtajzqxybqebd6+hy2emfmVGaC4DTiHcVZkcdRhB
K3GtsnOMrmQsi5rWFOHae4l405LGLn94BVt9+kkkJmBg7vDhBHIwxo4o1lE74g+b
Sf58cz9yzkatponn9f8/XvGGsweY2MV2ey/SMqmN9l+mjFMtkE65ZZ8p7w/RNEpR
uh3qHV9oMdFV/66C2kMev0SloNo/ZSgIhwIF63xX7FX94a+Bu9cmTWne9osp7mNY
KNV0pbDa2vcVk4VeTpml8vGK/aZ13wz7bbVJpwzk3o0VkAF8Cw7rvr9OnXRFRvaB
XAiLN26+joZlaGVcciEG3xzuaC0Ac1AefUd1dL+qWjeWEvNKrKfSI7vqpJfTJUh2
nHsHObNlFKLbFqGYZxM6jCnWdyGRHef1FDg=
-----END CERTIFICATE-----`

const selfSignedCAKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
