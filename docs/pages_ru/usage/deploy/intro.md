---
title: Введение
permalink: usage/deploy/intro.html
---

## Введение

За развертывание в werf отвечает Helm 3. werf стремится сделать работу с Helm более простой, удобной и гибкой, при этом не ломая обратную совместимость с Helm-чартами, Helm-шаблонами и Helm-релизами.

werf дает ряд дополнительных возможностей, недоступных пользователям оригинального Helm 3:

- умное отслеживание состояния ресурсов при развертывании;

- задание порядка развертывания для любых ресурсов, а не только для хуков;

- ожидание создания и готовности ресурсов, не принадлежащих релизу;

- интеграция сборки и развертывания и многое-многое другое.

## Развертывание с werf

Для развертывания простого приложения достаточно двух файлов и команды `werf converge`:

```
# .helm/templates/hello.yaml:
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello
spec:
  selector:
    matchLabels:
      app: hello
  template:
    metadata:
      labels:
        app: hello
    spec:
      containers:
      - image: nginxdemos/hello:plain-text
```

```
# werf.yaml:
configVersion: 1
project: hello
```

```shell
werf converge --repo registry.example.org/repo --env production
```

Результат: Deployment `hello` развёрнут в Namespace `hello-production`.

Более сложный пример развертывания со сборкой образов и внешними Helm-чартами:

```
# werf.yaml:
configVersion: 1
project: myapp
---
image: backend
dockerfile: Dockerfile
```

```dockerfile
# Dockerfile:
FROM node
WORKDIR /app
COPY . .
RUN npm ci
CMD ["node", "server.js"]
```

```yaml
# .helm/Chart.yaml:
dependencies:
- name: postgresql
  version: "~12.1.9"
  repository: https://charts.bitnami.com/bitnami
```

```yaml
# .helm/values.yaml:
backend:
  replicas: 1
```

{% raw %}

```
# .helm/templates/backend.yaml:
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
spec:
  replicas: {{ $.Values.backend.replicas }}
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - image: {{ $.Values.werf.image.backend }}
```

{% endraw %}

```shell
werf converge --repo registry.example.org/repo --env production
```

Результат: собран образ `backend`, а затем Deployment `backend` и ресурсы чарта `postgresql` развёрнуты в Namespace `myapp-production`.

Удалить развернутое приложение можно командой `werf dismiss`:

```shell
werf dismiss --repo registry.example.org/repo --env production
```