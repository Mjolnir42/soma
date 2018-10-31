# DESCRIPTION

This command is used to update existing users in SOMA. It is a full
replace, requiring the full user record to be specified in the command.
It is not possible to change the user account's system or activation
status flag with this command.

It is possible to mark user accounts as deleted using the update command
by including the optional `deleted true` keyword. It is not possible to
resurrect a deleted account by specifying `deleted false` on a deleted
account.

# SYNOPSIS

```
soma user-mgmt update ${userID} username ${uname} firstname ${fname} lastname ${lname} employeenr ${num} mailaddr ${addr} team ${team} [deleted ${bool}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
userID | string | UUID of the user | | no
uname | string | Username of the user | | no
fname | string | Given name of the user | | no
lname | string | Family name of the user | | no
num | integer | Numeric employee ID number of the user | | no
addr | string | Email address of the user | | no
team | string | Name of the team the user belongs to | | no
bool | boolean | Boolean flag if this user is deleted | false | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | identity | | no | yes
identity | user-mgmt | update | yes | no

# EXAMPLES

```
soma user-mgmt update ffffffff-ffff-ffff-ffff-ffffffffffff \
    username AnonymousCoward \
    firstname Anonymous \
    lastname Coward \
    employeenr 9999999999999999 \
    mailaddr devzero@example.com \
    team wheel \
soma user-mgmt update eed11901-73e1-46f9-8e98-12b582933eda \
    username jd \
    firstname Jon \
    lastname Doe \
    employeenr 1234 \
    mailaddr jon@example.com \
    team ExampleTeam \
    deleted true
```
