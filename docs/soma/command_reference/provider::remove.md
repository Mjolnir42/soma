# DESCRIPTION

This command is used to remove metric provider definitions from SOMA.

Provider names may not contain / characters.

# SYNOPSIS

```
soma provider remove ${provider}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
provider | string | Unit of the provider | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | provider | remove | yes | no

# EXAMPLES

```
soma provider remove collectd
```
