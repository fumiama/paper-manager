import type { MenuModule } from '/@/router/types'
import { t } from '/@/hooks/web/useI18n'
const menu: MenuModule = {
  orderNo: 20,
  menu: {
    name: t('routes.templist.name'),
    path: '/templist',

    children: [
      {
        path: 'index',
        name: t('routes.templist.name'),
      },
    ],
  },
}
export default menu
