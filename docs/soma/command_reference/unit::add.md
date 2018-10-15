# DESCRIPTION

This command is used to add unit definitions to SOMA.

The unit glyph string may not contain / characters.

# SYNOPSIS

```
soma unit add ${unit} name ${name}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
unit | string | Glyph of the unit | | no
name | string | Name of the unit | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | unit | add | yes | no

# EXAMPLES

```
soma unit add s name second
soma unit add b name bit
soma unit add B name Byte
```
