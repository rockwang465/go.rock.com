package license

//
//const defaultMasterCACert = `
//-----BEGIN CERTIFICATE-----
//MIIDozCCAougAwIBAgIJAJgOlpYuWlSUMA0GCSqGSIb3DQEBCwUAMGgxCzAJBgNV
//BAYTAkNOMRAwDgYDVQQIDAdCRUlKSU5HMRAwDgYDVQQHDAdCRUlKSU5HMRIwEAYD
//VQQKDAlTRU5TRVRJTUUxITAfBgNVBAMMGHByaXZhdGUuY2Euc2Vuc2V0aW1lLmNv
//bTAeFw0xNzEyMDYwOTM0MzdaFw0yNzEyMDQwOTM0MzdaMGgxCzAJBgNVBAYTAkNO
//MRAwDgYDVQQIDAdCRUlKSU5HMRAwDgYDVQQHDAdCRUlKSU5HMRIwEAYDVQQKDAlT
//RU5TRVRJTUUxITAfBgNVBAMMGHByaXZhdGUuY2Euc2Vuc2V0aW1lLmNvbTCCASIw
//DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALPbG4PxtqX9TEk720hkxqlY07WB
//KWg3MD51jzZzVEDe0LnsD0kmdSt0lA+WvIGwXNXh0TNX9B7zcNwJ+dhj6oEujA+Z
//zmd3FpulpJElU0nE/R68LzTa/4bXCIwMmpkKvMbuLdwSNimbSKiO9IGrloCNFTfP
//Fskmmp3NbcXkNFQCRseGFUGGJDfsNdSp5qGsTIolpqoBRlHyxsHxqzk3PVkvRZ0u
//7ytQKQENbb4w60ukqh45hLX6J0irQfqSY8Bw51gos3OfQ3ur8z3HdFMp+/PxMh4n
//rAMvqBLe4d6fBj+oj2Ej27gQZ8aDvV1jWh92rN5A9RKTM3XV90PRGHzMvn0CAwEA
//AaNQME4wHQYDVR0OBBYEFBc2fH74sxyPX/N+TbATRDVmcM1+MB8GA1UdIwQYMBaA
//FBc2fH74sxyPX/N+TbATRDVmcM1+MAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEL
//BQADggEBADScq9hKnAFlGw5gWJoNuTx6FPD2MJ6Zm0/VoD7xNS32nIaVVI0Tt6VH
//eZe0JD7Cer4LIPUb5oJTmcR2mUYgBhVLtZKoLRwgH7daRqaI/LOdV8XQR+qRqyj6
//iBtYOZmumXqvsW2NsrxV/fAWbXeVZl3bE7YVfbvktBhdFNT05DVEJDu+0QmoClHN
//e39TYZbLuUgfBIVZUVItKJfp1NVVX6M5U+/KEzwxShAVOez/S3Jsn+dROKBf6WQn
//mLmCh5WMppaIbSjWatz2hBcqarh12gGQgNwyd+zyWbqtCddEdaxNW8WLj1Y8JLxH
//rO2hAGzKct7qiBd6mDCBJfSWIVxKU0Q=
//-----END CERTIFICATE-----
//`
//
//const defaultSlaveCACert = `
//-----BEGIN CERTIFICATE-----
//MIIDrzCCApegAwIBAgIJAOI2xfBCEdAmMA0GCSqGSIb3DQEBCwUAMG4xCzAJBgNV
//BAYTAkNOMRAwDgYDVQQIDAdCRUlKSU5HMRAwDgYDVQQHDAdCRUlKSU5HMRIwEAYD
//VQQKDAlTRU5TRVRJTUUxJzAlBgNVBAMMHnNsYXZlLnByaXZhdGUuY2Euc2Vuc2V0
//aW1lLmNvbTAeFw0xNzEyMTIxMTA3MjNaFw0yNzEyMTAxMTA3MjNaMG4xCzAJBgNV
//BAYTAkNOMRAwDgYDVQQIDAdCRUlKSU5HMRAwDgYDVQQHDAdCRUlKSU5HMRIwEAYD
//VQQKDAlTRU5TRVRJTUUxJzAlBgNVBAMMHnNsYXZlLnByaXZhdGUuY2Euc2Vuc2V0
//aW1lLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALPbG4PxtqX9
//TEk720hkxqlY07WBKWg3MD51jzZzVEDe0LnsD0kmdSt0lA+WvIGwXNXh0TNX9B7z
//cNwJ+dhj6oEujA+Zzmd3FpulpJElU0nE/R68LzTa/4bXCIwMmpkKvMbuLdwSNimb
//SKiO9IGrloCNFTfPFskmmp3NbcXkNFQCRseGFUGGJDfsNdSp5qGsTIolpqoBRlHy
//xsHxqzk3PVkvRZ0u7ytQKQENbb4w60ukqh45hLX6J0irQfqSY8Bw51gos3OfQ3ur
//8z3HdFMp+/PxMh4nrAMvqBLe4d6fBj+oj2Ej27gQZ8aDvV1jWh92rN5A9RKTM3XV
//90PRGHzMvn0CAwEAAaNQME4wHQYDVR0OBBYEFBc2fH74sxyPX/N+TbATRDVmcM1+
//MB8GA1UdIwQYMBaAFBc2fH74sxyPX/N+TbATRDVmcM1+MAwGA1UdEwQFMAMBAf8w
//DQYJKoZIhvcNAQELBQADggEBAG8vG7uYYFpgwU6ZG1tVxjhMhMFnI7iIasX6kFrd
//7yi8N5T3PnYQfHY2ryCkZK6lkdqOhYjX7QuIptRhKeZtIKzkJZIzC2ImnQImf+ah
//WIkhN5pmuaA9rb43NRxnfCwLKbMxnheZnBUnFg/Ty83yYTcDEs2zAjNmiGJKLERn
//xIUnoWEiXb/tGTatTPNwmNtWbrfy3AeFP39iRD82FPXtsMve45+EnGpt2WAXjx/q
//LSbFMBojo7wGfFUu8rw7RDt9b8XgOgjQNYLUlct4MtIsCFMZJU17gCBJ5DFRTnHC
//MFD+L3DkGdtm5sbsgdsVB9F3vhsnFWO8y9E2uusM4G8rnT8=
//-----END CERTIFICATE-----
//`
//
//const defaultConsoleKey = `
//-----BEGIN RSA PRIVATE KEY-----
//MIIEowIBAAKCAQEAxV8OzmENKTgpshVjko98tT8aLeD61g+ujdVj9aviOKKm5Edh
//04jPwzr3n41ZMDf+B/xPqe7HWTxpv5Lu4Kqa1JBi0qqZZFShBqxcLQuAjVzNIGZe
//sDhIrgm93ubLP8ZwvK+vdZmPbOnx1VKfjeOglZKWS+VrwXG61IM+0iYPTXaOJX/J
//QDucSvcXeyfUxnybC4lgQdCQeTfF+nTWunO0a7A9vrbx79uzN5yUy+c5RBNuJhXd
//MqXDfWvI46hLqPJS2zTDnG5CVKBWm1N0G3HGtXBa4SiwgAn2g3My8TluIQU85ThQ
//Sy0umI//yMy8kXY5lLqxA2n72zfezS7PrHZ8zwIDAQABAoIBAAE7U6NUFbnxIMl8
//uq9ad+PFrgslQUt+s48tCr+ov/OsiDAahfDFBM7qGkuDnU/guZQhLfoYhGP5LYvF
//hfoe9nJnKEa6S9TFdm/NOZIKZVX8g0c1fFfLMiDr7KRsek4+lcuHqSepuqxqVVkI
//d/hxuDnWvVth5idB53GWFBlJpYTNOshNfllkN0+Gwyo5QZRt2aWTnyp8+g0wPq8t
//cVI1U5YB6tKABpl7qA25OhZhdHZF0tmwN51rsJ1YoML/PRpu1p3DH5FhJ81IDe0A
//U5lpBGkMuSORNftgba/8LfRwyqTn/YYn579aZl6579C4K1ANxWJJnO64zZahDCNV
//YrFz4WkCgYEA/drKBwlx+DPCSnX8f0thXNkrq9yk1l1m5+EIcd1oiHPl3iFHRJoo
//iHXh5P8RPvVxkzB9LPnnyEx9aMJ6tEONDYkf13K88YcR4JwMcVDdKLSGrpHK0mDc
//5xlkd/ueYwfiEKqnCwGOXw4BkPKopL5fyYaHLEn8GwyI1Mcgjw41KTMCgYEAxwoR
//XELzw4w0/nnzwBVwUfV76OdFrXylwktq1H7AMMQmfWXTIrBMLkY/f9mkALydZZeq
//pm8BjuKKKUicUN98Zh5LK7EQV2ogK3ps0OH6wQNVNnSz5oRRUsEkrM0Uaw15R613
//qyZ6Dg3OkDmr5ZyUmVROS45oDPoZlz4PfM63tfUCgYBk8DlCwQuzQIlx6CZFS2jk
//bWoDBVH59tuzOfSMqhglocf2Ik9fRNj3IcB3uMBXw2qstywe1SPHrjpzjFkUEoQk
//rLCfj3z3oNiH8iS0bg3yYI3pHgmCy4cq0Rr05nUdNYY7UE/pfW3p9/zBcOuDzjry
//O+7FuolnC/3gdWlJ2MFkpwKBgDALCxO1CXfbAPOn5iEoS5tM4OLf6B6vJqeWYqv2
//CFf9ELlV+be2zDyjMjKfCwoufOOHz2YrBzpBDk5Wu3x95V4U09ow/BvNfwRfoaJt
//2YP7VPc3BjGPIL4T5tFbEyGf9/VINsl2GSIJTSHc+dQLjobQJbHxJsZzG/g4v65F
//i2x9AoGBANYmI4o2ZGoa1ywxwVoOvQtl3JmVEwy8Wwfp9qv/9bmijjlkKBQFQNWC
//udksam8veH8inhBWok1G84T7e46IKcxSVz+26KaG3ViBIyVgrlJkSZ5JVBkZFRaG
//JyYxeOfzwryGz1Z6plgVStXq6O9RPuGUeVkT1K5CmeMM0PQlmaG0
//-----END RSA PRIVATE KEY-----
//`
//
//const defaultConsoleCert = `
//-----BEGIN CERTIFICATE-----
//MIIDdDCCAlwCCQD6c+kzCWZJDTANBgkqhkiG9w0BAQsFADB8MQswCQYDVQQGEwJD
//TjEQMA4GA1UECAwHQmVpamluZzEQMA4GA1UEBwwHQmVpamluZzESMBAGA1UECgwJ
//U2Vuc2V0aW1lMRIwEAYDVQQLDAlTZW5zZXRpbWUxITAfBgNVBAMMGHByaXZhdGUu
//Y2Euc2Vuc2V0aW1lLmNvbTAeFw0xODA3MTcwODU3MTdaFw0yODA3MTQwODU3MTda
//MHwxCzAJBgNVBAYTAkNOMRAwDgYDVQQIDAdCZWlqaW5nMRAwDgYDVQQHDAdCZWlq
//aW5nMRIwEAYDVQQKDAlTZW5zZXRpbWUxEjAQBgNVBAsMCVNlbnNldGltZTEhMB8G
//A1UEAwwYcHJpdmF0ZS5jYS5zZW5zZXRpbWUuY29tMIIBIjANBgkqhkiG9w0BAQEF
//AAOCAQ8AMIIBCgKCAQEAxV8OzmENKTgpshVjko98tT8aLeD61g+ujdVj9aviOKKm
//5Edh04jPwzr3n41ZMDf+B/xPqe7HWTxpv5Lu4Kqa1JBi0qqZZFShBqxcLQuAjVzN
//IGZesDhIrgm93ubLP8ZwvK+vdZmPbOnx1VKfjeOglZKWS+VrwXG61IM+0iYPTXaO
//JX/JQDucSvcXeyfUxnybC4lgQdCQeTfF+nTWunO0a7A9vrbx79uzN5yUy+c5RBNu
//JhXdMqXDfWvI46hLqPJS2zTDnG5CVKBWm1N0G3HGtXBa4SiwgAn2g3My8TluIQU8
//5ThQSy0umI//yMy8kXY5lLqxA2n72zfezS7PrHZ8zwIDAQABMA0GCSqGSIb3DQEB
//CwUAA4IBAQDFPS+zoXkDOAy6Y7dI0kwjpQlqZhKjPni1LAuCbNbebpGve9RlZQPr
//p0fu1zD8vwgnAwqOCJSPtdw2SphRQCmwOjkEazfHZezve3eJ3hIaMXO89Kwn14ye
//SGg9l/1/cwCun61kiQzvW2tMK8KxUtvO62TLmKZMmu6iMj1Koi98TsBnDHhMfpv1
//c1UTgXKLq+W7vJrMAwlQqbGI6xjxGG4AHpEkrK23qDlGqkJ1uNXtNf5+dQe14+j4
//MopZ2DjS2c2+Z0GfWW6D9IxoZVqq0eBdfLskc6THEh/JwZWbpIYL7k/zQELiB6HO
//UUuPtAJTIpTsOPk8nMmG93jLTQYERIa6
//-----END CERTIFICATE-----
//`
//
//const (
//	// Master presents Master CA
//	Master ServerType = iota
//	// Slave presents Slave CA
//	Slave
//)
//
//// ServerType defines master or slave server type
//type ServerType uint