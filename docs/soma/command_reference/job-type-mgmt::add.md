# DESCRIPTION

This command is used to add job type definitions to SOMA.

# SYNOPSIS

```
soma job type-mgmt add ${type}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
type | string | Name of the type | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | job-type-mgmt | add | yes | no

# EXAMPLES

```
soma job type-mgmt add repository::destroy
soma job type-mgmt add check-config::create
```
