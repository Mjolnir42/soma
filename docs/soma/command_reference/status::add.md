# DESCRIPTION

This command is used to add workflow status definitions to SOMA.

# SYNOPSIS

```
soma status add ${status}
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
global | status | add | yes | no

# EXAMPLES

```
soma status add rollout_in_progress
soma status add rollout_failed
```
