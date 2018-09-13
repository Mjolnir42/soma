# DESCRIPTION

This command is used to remove entity definitions from SOMA.

# SYNOPSIS

```
soma entity remove ${entity}
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
global | entity | remove | yes | no

# EXAMPLES

```
soma entity remove repository
soma entity remove bucket
```

