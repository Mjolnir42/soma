# DESCRIPTION

This command is used to remove users from SOMA. User accounts are
flagged as deleted but not removed from the database.

Once a user is removed, it can not be resurrected.

# SYNOPSIS

```
soma user-mgmt remove ${uname}
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
identity | user-mgmt | remove | yes | no

# EXAMPLES

```
soma user-mgmt remove jd
```
