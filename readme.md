# 2FA Example Service

#### Steps

1. Using a browser go to http://localhost:8080/QrCode

1. Using a OTP authentication app, scan the produced QrCode

1. Using the below command and your access code from the authenticator app, run the below `curl` command

```sh
curl -v -u John:IamNotACat localhost:8080/login --data '{"OTP":"ACCESS_CODE"}'
```

