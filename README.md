# update-yaml-values
This won't be helpful to anyone but me


pipeline-2.yml
```
resource_types:
  - name: gcs-resource
    type: registry-image
    source:
      repository: gcr.io/tas-ppe/pivotalcfreleng/gcs-resource
      username: _json_key
      password: ((dependabot.ppe_gcr_service_account_key))
resources: abc
```

Usage example

```
nr $ ./update-secrets ./fixtures/pipeline-2.yml
resource_types:
- name: gcs-resource
  source:
    password: ((dependabot.dependabot.ppe_gcr_service_account_key))
    repository: gcr.io/tas-ppe/pivotalcfreleng/gcs-resource
    username: _json_key
  type: registry-image
resources: abc

```
