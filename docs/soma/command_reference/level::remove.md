# DESCRIPTION

This command is used to remove notification levels to SOMA.
The level to remove may be specified by either its long or short
name.

Level names must not contain `/` characters.

# SYNOPSIS

```
soma level remove ${lvl}|${abbrev}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
lvl | string | Name of the level | | no
abbrev | string | Abbreviation of the level name | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | level | remove | yes | no

# EXAMPLES

```
soma level remove informational
soma level remove info 
```
