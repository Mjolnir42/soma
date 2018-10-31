# DESCRIPTION

This command is used to delete removed users from the SOMA database.

This action is only possible if the user account is not referenced
in any of the history or job tables.

# SYNOPSIS

```
soma user-mgmt purge ${uname}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
uname | string | Username of the user | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | identity | | no | yes
identity | user-mgmt | purge | yes | no

# EXAMPLES

```
soma user-mgmt purge jd
```
