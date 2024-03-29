import type { AppRouteModule } from '/@/router/types'

import { LAYOUT } from '/@/router/constant'
import { t } from '/@/hooks/web/useI18n'
import { RoleEnum } from '/@/enums/roleEnum'

const dashboard: AppRouteModule = {
  path: '/dashboard',
  name: 'Dashboard',
  component: LAYOUT,
  redirect: '/dashboard/analysis',
  meta: {
    orderNo: 10,
    icon: 'ion:grid-outline',
    title: t('routes.dashboard.dashboard'),
  },
  children: [
    {
      path: 'analysis',
      name: 'Analysis',
      component: () => import('/@/views/dashboard/analysis/index.vue'),
      meta: {
        // affix: true,
        title: t('routes.dashboard.analysis'),
        roles: [RoleEnum.SUPER],
      },
    },
    {
      path: 'workbench',
      name: 'Workbench',
      component: () => import('/@/views/dashboard/workbench/index.vue'),
      meta: {
        title: t('routes.dashboard.workbench'),
      },
    },
    {
      path: 'account',
      name: 'Account',
      component: () => import('/@/views/dashboard/account/index.vue'),
      meta: {
        // affix: true,
        title: t('routes.dashboard.account'),
        roles: [RoleEnum.SUPER],
      },
    },
    {
      path: 'regex',
      name: 'Regex',
      component: () => import('/@/views/dashboard/regex/index.vue'),
      meta: {
        // affix: true,
        title: t('routes.dashboard.regex'),
      },
    },
  ],
}

export default dashboard
