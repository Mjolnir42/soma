# DESCRIPTION

This command is used to delete a permission section from the system.

# SYNOPSIS

```
soma sections remove ${section} [from ${category}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
section | string | Name of the section | | no
category | string | Name of the category | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | permission | | no | yes
permission | section | remove | yes | no

# EXAMPLES

```
soma section remove environment
soma section remove datacenter from global
```
