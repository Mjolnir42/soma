# DESCRIPTION

This command is used to add metric provider definitions to SOMA.

Provider names may not contain / characters.

# SYNOPSIS

```
soma provider add ${provider}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
provider | string | Name of the provider | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | provider | add | yes | no

# EXAMPLES

```
soma provider add prometheus
soma provider add snap
```
