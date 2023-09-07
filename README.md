# self-signed-certificate
A script to generate self-signed certificate.

No need to modify any files.

The certificate generation directory is the out directory under the current directory. 
The root certificate is only generated once by default, if it exists. 
You can manually delete the directory or the root certificate under the directory to regenerate.

### Usage

Generate all certificates

```
./gen.cert.sh -h

Usage: ./gen.cert.sh [-d domain, optional parameter, multiple parameters separated by commas ] 
       [-i ip, optional parameter, multiple parameters separated by commas ] 
       [-s ca subject, defalut: /C=CN/ST=Jiangsu/L=Wuxi/O=zzz/OU=zzz/CN=zzz ] 
       [-D ca validity days, default: 73000 ] 
       [-rs root ca subject, If not specified same to -s ] 
       [-rD root ca validity days, If not specified same to -D ] 
       [-sn serial number, default: 1000 ]

example:
./gen.cert.sh -d 'test.com,*.zzz.com' -i '1.1.1.1,192.168.102.60'
```

Generate root certificate only

```
./gen.root.sh -h

Usage: ./gen.root.sh [-s root ca subject, defalut: /C=CN/ST=Jiangsu/L=Wuxi/O=zzz/OU=zzz/CN=zzz Root CA ] 
         [-d validity days, default: 73000 ] 
         [-sn serial number, default: 1000 ]
```

Delete certificate directory

```
rm -rf ./out/
```

### refer to

[Fishdrowned/ssl](https://github.com/Fishdrowned/ssl/)
