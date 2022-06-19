# DeBotOps

### 🤖 De(v) Bot Ops

DevOps automation system for kubernetes

## ⚙️ Requirements

- Kubernetes
    - Istio ( + Ingress Gateway )

## ⭐️ Features

- Application
    - `Application` describes application for service.

- Listener
    - `Listener` configure inbound traffics.

- Mapping
    - `Mapping` connect between `Application` and `Listener`.

## ✅ TODO

TODO list for first release.

- `Application`
    - [x] CRD
    - [x] Controller
    - [ ] API
    - [ ] WebHook

- `Listener`
    - [x] CRD
    - [x] Controller
    - [ ] API
    - [ ] WebHook

- `Mapping`
    - [x] CRD
    - [x] Controller
    - [ ] API
    - [ ] WebHook

- `Injection`
    - [ ] CRD
    - [ ] Controller
    - [ ] API
    - [ ] WebHook

- Integration
    - Network
        - [x] Istio
        - [ ] cert-manager

    - CI / CD
        - [ ] ArgoCD
