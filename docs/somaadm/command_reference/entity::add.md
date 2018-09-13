# DESCRIPTION

This command is used to add entity definitions to SOMA.

# SYNOPSIS

```
soma entity add ${entity}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
entity | string | Name of the entity | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | entity | add | yes | no

# EXAMPLES

```
soma entity add repository
soma entity add bucket
```
