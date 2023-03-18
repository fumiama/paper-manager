import type { AppRouteModule } from '/@/router/types'

import { LAYOUT } from '/@/router/constant'
import { t } from '/@/hooks/web/useI18n'

const settings: AppRouteModule = {
  path: '/settings',
  name: 'Settings',
  component: LAYOUT,
  redirect: '/settings/index',
  meta: {
    hideChildrenInMenu: true,
    icon: 'ion:settings-outline',
    title: t('routes.settings.name'),
    orderNo: 200,
  },
  children: [
    {
      path: 'index',
      name: 'SettingsPage',
      component: () => import('/@/views/page/settings/index.vue'),
      meta: {
        title: t('routes.settings.name'),
        icon: 'ion:settings-outline',
        hideMenu: true,
      },
    },
    {
      path: 'password',
      name: 'PasswordSettingsPage',
      component: () => import('/@/views/page/settings/password/index.vue'),
      meta: {
        title: t('routes.settings.password'),
        icon: 'ion:settings-outline',
        hideMenu: true,
      },
    },
  ],
}

export default settings
