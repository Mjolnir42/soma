# somaadm permission show

This command is used to show details about a specific permission.

The permission to show can be specified via shorthand or regular
syntax.

# SYNOPSIS

```
somaadm permission show ${permission} in ${category}
somaadm permission show ${category}::${permission}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
permission | string | Name of the permission | | no
category | string | Name of the category | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | permission | | no | yes
permission | permission | show | yes | no

# EXAMPLES

```
./somaadm permission show browse in global
./somaadm permission show permission::auditor
```
