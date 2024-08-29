# 牙木ExternalDNS Webhook插件

[ExternalDNS](https://github.com/kubernetes-sigs/external-dns) 自动发现k8s对外发布的服务并将服务注册到外部DNS服务中。牙木提供ExternalDNS插件帮助客户将企业容器云平台发布的服务自动注册到牙木SmartDDI设备中并提供解析。

## 支持以下资源记录
| 记录类型    | 是否支持     |
|-------------|--------------|
| A           | 支持         |
| AAAA        | 支持         |
| CNAME       | 支持         |
| TXT         | 不支持       |
| PTR         | 不支持       |



## 版本要求

- ExternalDNS >= v0.14.0
- SmartDDI >= 3.6.2

## 部署

1. 在SmartDDI管理系统添加用户. `系统 > 用户管理 > 用户`

2. 给创建好的用户赋予第三方API权限以及权威管理功能权限

3. 添加external-dns Helm仓库

    ```sh
    helm repo add external-dns https://kubernetes-sigs.github.io/external-dns/
    ```

5. 创建k8s密钥配置 `external-dns-yamu-secret`  `api_user` 以及 `api_secret`值要跟网管创建的用户一致:

    ```yaml
    apiVersion: v1
    stringData:
      api_user: <INSERT USER NAME>
      api_secret: <INSERT USER KEY>
    kind: Secret
    metadata:
      name: external-dns-yamu-secret
    type: Opaque
    ```

6. 创建helm文件
7. `external-dns-webhook-values.yaml`:

    ```yaml
    fullnameOverride: external-dns-yamu
    logLevel: debug
    provider:
      name: webhook
      webhook:
        image:
          repository: ghcr.io/yamu-oss/external-dns-yamu-webhook
          tag: main # replace with a versioned release tag
        env:
          - name: YAMU_API_USER
            valueFrom:
              secretKeyRef:
                name: external-dns-yamu-secret
                key: api_user
          - name: YAMU_API_KEY
            valueFrom:
              secretKeyRef:
                name: external-dns-yamu-secret
                key: api_secret
          - name: YAMU_HOST
            value: https://192.168.1.1 # 替换为SmartDDI网管地址
          - name: LOG_LEVEL
            value: debug
          - name: VIEW
            value: "default" # 替换为客户默认视图
          - name: DEFAULT_TTL
            value: "600" # 替换为客户默认TTL
        livenessProbe:
          httpGet:
            path: /healthz
            port: http-wh-metrics
          initialDelaySeconds: 10
          timeoutSeconds: 5
        readinessProbe:
          httpGet:
            path: /readyz
            port: http-wh-metrics
          initialDelaySeconds: 10
          timeoutSeconds: 5
    extraArgs:
      - --ignore-ingress-tls-spec
    policy: sync
    sources: ["ingress", "service", "crd"]
    registry: noop
    domainFilters: ["yamu.com"] # 替换为客户域名
    ```

7. 安装

    ```sh
    helm install external-dns-yamu external-dns/external-dns -f external-dns-yamu-values yaml --version 1.14.5 -n external-dns
    ```

