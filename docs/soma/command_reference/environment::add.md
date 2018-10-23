# DESCRIPTION

This command is used to add environment definitions to SOMA.

Environment names must not contain / characters.

# SYNOPSIS

```
soma environment add ${environment}
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
global | environment | add | yes | no

# EXAMPLES

```
soma environment add testing
soma environment add prod
```
