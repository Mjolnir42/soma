# DESCRIPTION

This command is used to show details about a notification level in SOMA.
The level to show may be specified by either its long or short
name.

Level names must not contain `/` characters.

# SYNOPSIS

```
soma level show ${lvl}|${abbrev}
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
global | level | show | yes | no

# EXAMPLES

```
soma level show error
soma level show err
```
