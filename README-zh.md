# Let me in
自动添加当前外网 IP 到 AWS 防火墙白名单
[English](./README.md)

# 干啥的

将您的服务公开到公共领域非常危险。
特别是一些设计不佳的服务，或者这些可用于挖矿服务（K8，YARN等），非常容易被黑客抓肉鸡用来挖矿。
如果必须执行此操作，则过滤源IP地址是相对安全的方法。

# 使用

0. 编译

```bash
go build
```

1. 配置
   
```yaml
Auth:
  UserPass:
    auxten: "123456"              # 用户名&密码 for HTTP Basic Authentication
AwsSg:
  Region: "cn-northwest-1"        # AWS Region
  SgName: "Hadoop"                # Security Group Name
  SgId:   "sg-0e0c5cd076cf1fb51"  # Security Group I
```
2. 运行

```bash
# AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY 可以在 AWS 界面上的 "My security credentials" 获取
export AWS_ACCESS_KEY_ID=XXX AWS_SECRET_ACCESS_KEY=XXXX 
./letmein config.yaml
```

3. 开门

进入 http://host:1323/ping ，输入用户名和密码。
`letmein`将向安全组添加一个新规则，该规则将允许来自源IP的所有流量通过。
因此，您应该在“安全组”内的主机上运行`letmein`。
