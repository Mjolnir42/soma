# DESCRIPTION

This command is used to show job type definitions from SOMA.

# SYNOPSIS

```
soma job type-mgmt show ${type}
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
global | job-type-mgmt | show | yes | no

# EXAMPLES

```
soma job type-mgmt show bucket::rename
soma job type-mgmt show check-config::destroy
```
