# DESCRIPTION

This command is used to rename an entity definitions in SOMA. Outside of
typofixing during system setup this command is probably unused.

# SYNOPSIS

```
soma entity rename ${entity} to ${new-entity}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
entity | string | Old name of the entity | | no
new-entity | string | New name of the entity | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | entity | rename | yes | no

# EXAMPLES

```
soma entity rename rpeository to repository
```
