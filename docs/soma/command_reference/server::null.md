# DESCRIPTION

This command initializes the special null-Server within SOMA. This
server is used as the default server for nodes which do not have a
server specified.

# SYNOPSIS

```
soma server null datacenter ${locode}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
locode | string | UN/Locode of the datacenter | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | server | null | yes | no

# EXAMPLES

```
soma server null datacenter de.fra
```
