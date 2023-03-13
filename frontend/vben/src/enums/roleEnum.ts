export enum RoleEnum {
  // super admin, all permission granted
  SUPER = 'super',

  // can only create / delete account of normal users
  ACCOUNT_MANAGER = 'accmgr',

  // add / del files + USER's permission
  FILE_MANAGER = `filemgr`,

  // have the permission of using the normal application
  USER = 'user',
}
