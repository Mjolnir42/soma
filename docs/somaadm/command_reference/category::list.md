# somaadm category list

This command is used to list permission categories.

# SYNOPSIS

```
soma category list
```

# ARGUMENT TYPES

This command takes no arguments.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | permission | | no | yes
permission | category | list | yes | no

# EXAMPLES

```
soma category list
```
