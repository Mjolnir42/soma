# DESCRIPTION

This command is used to remove view definitions from SOMA.

# SYNOPSIS

```
soma view remove ${view}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
view | string | Name of the view | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | view | remove | yes | no

# EXAMPLES

```
soma view remove local
soma view remove any
```
