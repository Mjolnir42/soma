# DESCRIPTION

This command is used to remove an attribute. This is currently only
possible if the attribute is unused.

# SYNOPSIS

```
soma attribute remove ${attribute}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
attribute | string | Name of the attribute | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | attribute | remove | yes | no

# EXAMPLES

```
soma attribute remove app_proto
```
