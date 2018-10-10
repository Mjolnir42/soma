# DESCRIPTION

This command is used to add a metric definition to SOMA. The definition can
include package specifications, providing the information which package must
be installed to collect this metric depending on the specific metric
provider system.

Metric names may not contain / characters.

# SYNOPSIS

```
soma metric add ${metric} unit ${unit} description ${text} \
                          [package ${provider}::${package}, ...]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
metric | string | Name of the metric | | no
unit | string | Glyph of the unit | | no
text | string | Description of the metric | | no
package | string | Per-provider package specification | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | metric | add | yes | no

# EXAMPLES

```
soma metric add icmp.echo.rtt unit ms description 'ICMP echo request roundtrip time'
soma metric add process.count unit count description 'Number of processes' \
            package stats::process package collect::local
```
