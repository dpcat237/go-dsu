## OSS Index

To check vulnerabilities from OSS Index is necessary to pass the OSS Index API token.

Registration on OSS Index is free following [this link](https://ossindex.sonatype.org/user/signin). 
After registration, you can find API token on [settings page](https://ossindex.sonatype.org/user/settings).

To use `go-dsu` with OSS Index you can pass email and API token (Eg. `--ossemail=email --osstoken=token`) 
or encode email wth API token via base64 (Eg. `--oss=base64-token`).

Encode token on Unix `$ echo -n 'email:token' | base64`.