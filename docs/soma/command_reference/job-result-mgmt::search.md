# DESCRIPTION

This command is used to search job result definitions from SOMA
by id or name.

At least one search criteria must be specified. If both are specified,
then both must match.

# SYNOPSIS

```
soma job result-mgmt search [id ${uuid}] [name ${result}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
uuid | string | UUID of the result | | yes
result | string | Name of the result | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | job-result-mgmt | search | yes | no

# EXAMPLES

```
soma job result-mgmt search name pending
soma job result-mgmt search id 0eba9b56-ebc3-42f0-9cc0-97eacf7d9e26
soma job result-mgmt search id 466173c7-ed9d-4940-b8f3-1935aa0775f4 name success
```
