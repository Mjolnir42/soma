# DESCRIPTION

This command is used to add job result definitions to SOMA.

# SYNOPSIS

```
soma job result-mgmt add ${result}
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
global | job-result-mgmt | add | yes | no

# EXAMPLES

```
soma job result-mgmt add pending
soma job result-mgmt add success
soma job result-mgmt add failed
```
