# sshlistener
Opens a ssh server and listens for authentication data (login and password). Input is logged.

This code is created to collect most popular passwords via AWS based experiment: <a href="http://romanenco.com/collect-passwords" target="_blank">collect passwords blog post</a>.

### Generate SSH key (optional)
```
ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
```
Leave passphrase empty.

### Run listener
```
go run sshlistener.go -private=id_rsa -output=log.out -port=2022
```
Keys:

 - private: path to a private key
 - output: path to a log file to be created or appended
 - port: to listen at


### Connect to the listener
```
ssh testuser@192.168.0.14 -p2022
```
Listener always replies with invalid password, but data is logged into configured log file.
