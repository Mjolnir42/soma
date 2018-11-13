# DESCRIPTION

This command is used to search job type definitions from SOMA
by id or name.

At least one search criteria must be specified. If both are specified,
then both must match.

# SYNOPSIS

```
soma job type-mgmt search [id ${uuid}] [name ${type}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
uuid | string | UUID of the type | | yes
type | string | Name of the type | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | job-type-mgmt | search | yes | no

# EXAMPLES

```
soma job type-mgmt search name cluster::property-destroy
soma job type-mgmt search id 4059000a-cd3d-489b-b5ce-585c08978cf4
soma job type-mgmt search id 4059000a-cd3d-489b-b5ce-585c08978cf4 name cluster::property-destroy
```
