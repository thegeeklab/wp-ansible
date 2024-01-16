---
title: wp-ansible
---

[![Build Status](https://ci.thegeeklab.de/api/badges/thegeeklab/wp-ansible/status.svg)](https://ci.thegeeklab.de/repos/thegeeklab/wp-ansible)
[![Docker Hub](https://img.shields.io/badge/dockerhub-latest-blue.svg?logo=docker&logoColor=white)](https://hub.docker.com/r/thegeeklab/wp-ansible)
[![Quay.io](https://img.shields.io/badge/quay-latest-blue.svg?logo=docker&logoColor=white)](https://quay.io/repository/thegeeklab/wp-ansible)
[![Go Report Card](https://goreportcard.com/badge/github.com/thegeeklab/wp-ansible)](https://goreportcard.com/report/github.com/thegeeklab/wp-ansible)
[![GitHub contributors](https://img.shields.io/github/contributors/thegeeklab/wp-ansible)](https://github.com/thegeeklab/wp-ansible/graphs/contributors)
[![Source: GitHub](https://img.shields.io/badge/source-github-blue.svg?logo=github&logoColor=white)](https://github.com/thegeeklab/wp-ansible)
[![License: Apache-2.0](https://img.shields.io/github/license/thegeeklab/wp-ansible)](https://github.com/thegeeklab/wp-ansible/blob/main/LICENSE)

Woodpecker CI plugin to manage infrastructure with [Ansible](https://github.com/ansible/ansible).

<!-- prettier-ignore-start -->
<!-- spellchecker-disable -->
{{< toc >}}
<!-- spellchecker-enable -->
<!-- prettier-ignore-end -->

## Usage

```YAML
kind: pipeline
name: default

steps:
  - name: ansible
    image: quay.io/thegeeklab/wp-ansible
    settings:
      playbook: deployment/playbook.yml
      private_key:
        from_secret: ansible_private_key
      inventory: deployment/hosts.yml
```

### Parameters

<!-- prettier-ignore-start -->
<!-- spellchecker-disable -->
{{< propertylist name=wp-ansible.data sort=name >}}
<!-- spellchecker-enable -->
<!-- prettier-ignore-end -->

## Build

Build the binary with the following command:

```Shell
make build
```

Build the Container image with the following command:

```Shell
docker build --file Containerfile.multiarch --tag thegeeklab/wp-ansible .
```

## Test

```Shell

```
