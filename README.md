# Defectdojo Resource

Gives a way to push reports to [defectdojo](https://github.com/DefectDojo/django-DefectDojo).  For right now, `check` and `in` will be noops.

Future state?  `In` can be used to get the security posture of a given application and build maybe?  A way to "quality-gate" based on aggregated security scans.

## Source Configuration

* `defectdojo_url`: *Required.* The path of the hosted instance of defectdojo.

* `debug`: *Optional.* Enable debug logging.

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

## Behavior

### `out`

Pushes a report of a given type to Defectdojo for the specified application.

#### Parameters

* `report_type`: *Required.* The type of report you're trying to upload.  The format of this string must match the internal scan type strings that defectdojo is using, [found here](https://github.com/DefectDojo/django-DefectDojo/blob/b08723ded1491d82910e51810de27963ee6ccca2/dojo/tools/factory.py).