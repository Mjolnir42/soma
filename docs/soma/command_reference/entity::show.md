# DESCRIPTION

This command shows details for an entity defined in SOMA.

# SYNOPSIS

```
soma entity show ${entity}
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
global | entity | show | yes | no

# EXAMPLES

```
soma entity show bucket
soma entity show cluster
```
