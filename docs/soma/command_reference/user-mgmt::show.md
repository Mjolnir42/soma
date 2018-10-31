# DESCRIPTION

This command is used to show details about users in SOMA.

It is possible to show details about removed user accounts that have not
been purged yet.

# SYNOPSIS

```
soma user-mgmt show ${uname}
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
identity | user-mgmt | show | yes | no

# EXAMPLES

```
soma user-mgmt show jd
```
