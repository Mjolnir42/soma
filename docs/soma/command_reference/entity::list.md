# DESCRIPTION

This command lists all entities defined in SOMA.

# SYNOPSIS

```
soma entity list
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | entity | list | yes | no

# EXAMPLES

```
soma entity list
```

