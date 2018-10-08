# DESCRIPTION

This command is used to remove monitoring system mode definitions from SOMA.

# SYNOPSIS

```
soma mode remove ${entity}
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
global | mode | remove | yes | no

# EXAMPLES

```
soma mode remove shared
```
