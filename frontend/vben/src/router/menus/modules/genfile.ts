import type { MenuModule } from '/@/router/types'
import { t } from '/@/hooks/web/useI18n'
const menu: MenuModule = {
  orderNo: 40,
  menu: {
    name: t('routes.genfile.name'),
    path: '/genfile',

    children: [
      {
        path: 'index',
        name: t('routes.genfile.name'),
      },
    ],
  },
}
export default menu
