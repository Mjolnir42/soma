# DESCRIPTION

This command is used to remove workflow status definitions from SOMA.

# SYNOPSIS

```
soma status remove ${status}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
status | string | Name of the status | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | status | remove | yes | no

# EXAMPLES

```
soma status remove active
soma status remove awaiting_deletion
```
