import type { MenuModule } from '/@/router/types'
import { t } from '/@/hooks/web/useI18n'
const menu: MenuModule = {
  orderNo: 30,
  menu: {
    name: t('routes.templist.name'),
    path: '/templist',

    children: [
      {
        path: 'index',
        name: t('routes.templist.templist'),
      },
      {
        path: 'chkdup',
        name: t('routes.templist.dup'),
      },
    ],
  },
}
export default menu
