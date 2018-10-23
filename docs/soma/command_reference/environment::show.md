# DESCRIPTION

This command shows details for an environment defined in SOMA.

# SYNOPSIS

```
soma environment show ${environment}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
environment | string | Name of the environment | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | environment | show | yes | no

# EXAMPLES

```
soma environment show prelive
```
