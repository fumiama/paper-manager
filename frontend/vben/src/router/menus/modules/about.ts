import type { MenuModule } from '/@/router/types'
import { t } from '/@/hooks/web/useI18n'
const menu: MenuModule = {
  orderNo: 100,
  menu: {
    name: t('routes.dashboard.about'),
    path: '/about',

    children: [
      {
        path: 'index',
        name: t('routes.dashboard.about'),
      },
    ],
  },
}
export default menu
