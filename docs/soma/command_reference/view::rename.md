# DESCRIPTION

This command is used to rename an view definitions in SOMA. Outside of
typofixing during system setup this command is probably unused.

# SYNOPSIS

```
soma view rename ${view} to ${new-view}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
view | string | Old name of the view | | no
new-view | string | New name of the view | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | view | rename | yes | no

# EXAMPLES

```
soma view rename int to internal
```
