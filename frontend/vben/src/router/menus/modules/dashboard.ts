import type { MenuModule } from '/@/router/types'
import { t } from '/@/hooks/web/useI18n'
const menu: MenuModule = {
  orderNo: 10,
  menu: {
    name: t('routes.dashboard.dashboard'),
    path: '/dashboard',

    children: [
      {
        path: 'analysis',
        name: t('routes.dashboard.analysis'),
      },
      {
        path: 'workbench',
        name: t('routes.dashboard.workbench'),
      },
      {
        path: 'account',
        name: t('routes.dashboard.account'),
      },
      {
        path: 'regex',
        name: t('routes.dashboard.regex'),
      },
    ],
  },
}
export default menu
