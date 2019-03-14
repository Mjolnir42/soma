# DESCRIPTION

This command is used to list all grants of a specific permission.

# SYNOPSIS

```
soma right list ${category}::${permission}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
category | string | Name of the category | | no
permission | string | Name of the permission | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | permission | | no | yes
permission | right | list | yes | no

# EXAMPLES

```
soma right list global::browse
```
