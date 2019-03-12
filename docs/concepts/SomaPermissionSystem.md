# SOMA Permission System

## Object Types

The SOMA permission model relies on five object types.

1. `action`, correspond to request executions by the application
2. `section`, groups related actions that are performed by the same
   application handler or on the same object
3. `permission`, references the sections and actions that it permits
4. `category`, by-scope groups of sections and permissions
5. `grant`, references the permission it grants as well as the user it
   was granted to

## Notation

`category`, `section`, `action` and `permission` can not be created
containing the `:` symbol. Full notation for an action is
`category::section:action`, full notation for a permission is
`category::permission`.

## Object Relations

The relations between these objects are defined as follows.

1. `action` belongs to exactly one `section`
2. `section` belongs to exactly one `category`
3. `permission` belongs to exactly one `category`
4. `action` can be mapped to a `permission` of the same `category`
5. `section` can be mapped to a `permission` of the same `category`,
    this represents all `action` of that `section` 
6. multiple `action` and `section` can be mapped to a `permission`
7. one `action` or `section` can be mapped to multiple `permission`
8. a user can be granted multiple `permission`

## Grant categories and permissions

Grant categories allow to distribute the permission management to users.

1.  Creating a category `foobar` also creates `foobar:grant`.
2.  Creating a permission `barfoo` in `foobar` also creates `barfoo` in
    `foobar:grant`. `foobar:grant::barfoo` is linked as granting
    `foobar::barfoo`.
3.  When a user has a grant permission, that user can give that
    permission to other users. For scoped permission, if a user has a
    grant permission on a specific object, the user can give out that
    permission on that object.

For example a user that has been granted the permission
`monitoring:grant::monitoringsystem:use` on a specific monitoring system can
grant other users the ability to deploy checks on that monitoring system.

## Implementation vs runtime definition

1.  Actions are implementation defined as they are used by request handlers
    to query authorization
2.  Sections are implementation defined as sections group actions, and
    roughly correspond to SOMA application handlers.
3.  Categories are implementation defined as they represent the permission
    scopes.
4.  Permissions are runtime defines.
5.  Mapping actions and sections to permissions are runtime defines.

## Magic

`somadbctl` installs the following during database schema installation:

1.  category `omnipotence`, ID `00000000-0000-0000-0000-000000000000`
2.  permission `omnipotence`, ID `00000000-0000-0000-0000-000000000000`
3.  group `wheel`, ID `00000000-0000-0000-0000-000000000000`
4.  user `root, ID `00000000-0000-0000-0000-000000000000`
5.  user `root` is granted permission `omnipotence`
6.  category `system` with a random ID
7.  a random bearer token that allows to authenticate and initialize the
    root user

### Omnipotence

1.  only root can have `omnipotence`, enforced via database check constraints
2.  `omnipotence` authorizes every request for every action
3.  `omnipotence` can not be granted, you either have it or you don't

### System

1.  Category `system` is installed during the database setup.
2.  There is no `system:grant` category, to grant `system` permissions one
    must have `omnipotence` or `system::permission`
3.  It is not possible to create permissions in category `system`
4.  Whenever a new category is created, a permission is created in
    `system` with the category name as the permission name.
5.  A system permission grants access to all actions in the category
    with the same name.
6.  System permissions are completely unscoped and allow all of their
    category's actions across all scopes.
7.  System permissions can only be granted to admin accounts or tool accounts
    with active isSystem flag
8.  admins and system tools can not have non-system permissions

# Appendix

## Category List

1.  `omnipotence`
2.  `system`
3.  `global`, unscoped permissions
4.  `identity`, unscoped identity management related permissions
6.  `monitoring`, per monitoring system permissions
8.  `operation`, unscoped volatile system administration permissions
7.  `permission`, unscoped permission related permissions
4.  `repository`, per repository permissions
5.  `self`,  oneself-scoped permissions
5.  `team`,  per team permissions
