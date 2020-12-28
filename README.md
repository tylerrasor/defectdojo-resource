# Defectdojo Resource

Gives a way to push reports to [defectdojo](https://github.com/DefectDojo/django-DefectDojo).  For right now, `check` and `in` will be noops.

Future state?  `In` can be used to get the security posture of a given application and build maybe?  A way to "quality-gate" based on aggregated security scans.

## Source Configuration

* `defectdojo_url`: *Required.* The path of the hosted instance of defectdojo.

### Example

``` yaml
resource_types:
- name: defectdojo-resource
  type: registry-image
  source:
    repository: tylerrasor/defectdojo-resource
    tag: latest

resources:
- name: defectdojo
  type: defectdojo-resource
  source:
    defectdojo_url: https://path-to-your-hosted-instance.io
```