# Defectdojo Resource

Gives a way to push reports to [defectdojo](https://github.com/DefectDojo/django-DefectDojo).  For right now, `check` and `in` will be noops.

Future state?  `In` can be used to get the security posture of a given application and build maybe?  A way to "quality-gate" based on aggregated security scans.

## Source Configuration

| Parameter                     | Type   | Required | Default | Description |
|:------------------------------|:-------|:---------|:--------|:------------|
| `defectdojo_url`              | URL    | yes      |         | The path of the hosted instance of defectdojo. |
| `api_key`                     | string | yes      |         | Generated API key (for `username`) to interact with defectdojo, [see here](https://defectdojo.readthedocs.io/en/latest/api-v2-docs.html). |
| `product_name`                | string | yes      |         | Name of the product (application) in defectdojo that we want to interact with. |
| `debug`                       | bool   | no       | `false` | Enable debug logging. |

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
    api_key: ((from-a-secret-manager-probably))
    app_name: "myApp"

jobs:
- name: scan-and-report
  plan:
  - task: do-the-scan
    config:
      platform: linux
      image_resource:
        type: registry_image
        source: { repository: alpine }
      output:
      - name: reports
      run:
        path: sh
        args:
          - |
            echo "do some cool scan" > reports.txt
  - put: defectdojo
    params:
      report_type: demo-scan
      path_to_report: reports/report.txt
```

## Behavior

### `out`

Pushes a report of a given type to Defectdojo for the specified application.

#### Parameters

| Parameter        | Type   | Required | Default | Description |
|:-----------------|:-------|:---------|:--------|:------------|
| `report_type`    | string | yes      |         | The type of report you're trying to upload.  The format of this string must match the internal scan type strings that defectdojo is using, [found here](https://github.com/DefectDojo/django-DefectDojo/blob/b08723ded1491d82910e51810de27963ee6ccca2/dojo/tools/factory.py). |
| `path_to_report` | string | yes      |         | File path (passed in from previous task) to the report you're trying to upload. |
| `make_active`    | bool   | no       | `false` | Should the scan be marked as `active` as far as findings go. |
