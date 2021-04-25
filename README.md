# Let me in
Authorize AWS Security Group Ingress with My Current IP(let me in)

[中文说明](./README-zh.md)

# Why
Exposing your service to public domain is very dangerous. 
Especially some poorly designed services, or these can be used to mining (K8s, YARN, etc). 
If you have to do this, filtering the SRC IP address is relatively a safe way.

# Usage

1. Config
   
```yaml
Auth:
  UserPass:
    auxten: "123456"              # Username & password for HTTP Basic Authentication
AwsSg:
  Region: "cn-northwest-1"        # AWS Region
  SgName: "Hadoop"                # Security Group Name
  SgId:   "sg-0e0c5cd076cf1fb51"  # Security Group I
```
2. Run

```bash
# AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY can be fetched from the "My security credentials"
export AWS_ACCESS_KEY_ID=XXX AWS_SECRET_ACCESS_KEY=XXXX 
./letmein config.yaml
```

3. Turn key

Access `http://the-server:1323/ping`, enter the username and password.
The `letmein` will add a new rule to the Security Group that let all traffic from the source IP pass.
So you should run `letmein` on the Host that inside the `Security Group`.

4. Revoke key

Access `http://the-server:1323/revoke/:ip`, ip is the string which you want revoke from aws sg.
The `letmein` will delete a rule from the Security Group.
