# DESCRIPTION

This command purges a removed server from SOMA by removing it from the
database.

For this command to work, the server must no longer be referenced.

# SYNOPSIS

```
soma server purge ${name}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
name | string | Name of the server | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | server | purge | yes | no

# EXAMPLES

```
soma server purge example-server-a
soma server purge example-server-b
```
