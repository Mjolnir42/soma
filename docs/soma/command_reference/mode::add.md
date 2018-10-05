# DESCRIPTION

This command is used to add monitoring system mode definitions to SOMA.

# SYNOPSIS

```
soma mode add ${mode}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
mode | string | Name of the mode | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | mode | add | yes | no

# EXAMPLES

```
soma mode add public
```
