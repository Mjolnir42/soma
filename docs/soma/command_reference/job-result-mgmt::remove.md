# DESCRIPTION

This command is used to remove job result definitions from SOMA.

# SYNOPSIS

```
soma job result-mgmt remove ${result}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
result | string | Name of the result | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | job-result-mgmt | remove | yes | no

# EXAMPLES

```
soma job result-mgmt remove pending
```
