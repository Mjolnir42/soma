# DESCRIPTION

This command is used to add view definitions to SOMA.

# SYNOPSIS

```
soma view add ${view}
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
global | view | add | yes | no

# EXAMPLES

```
soma view add local
soma view add any
```
