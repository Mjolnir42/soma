# DESCRIPTION

This command shows details for a workflow status defined in SOMA.

# SYNOPSIS

```
soma status show ${status}
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
global | status | show | yes | no

# EXAMPLES

```
soma status show deprovisioned
```
