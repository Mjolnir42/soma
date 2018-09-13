# DESCRIPTION

This command is used to show details about a specific permission.

The permission to show can be specified via shorthand or regular
syntax.

# SYNOPSIS

```
soma permission show [${category}::]${permission} [in ${category}]
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
soma permission show browse in global
soma permission show permission::auditor
```
