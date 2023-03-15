import type { MenuModule } from '/@/router/types'
import { t } from '/@/hooks/web/useI18n'
const menu: MenuModule = {
  orderNo: 20,
  menu: {
    name: t('routes.filelist.name'),
    path: '/filelist',

    children: [
      {
        path: 'index',
        name: t('routes.filelist.name'),
      },
    ],
  },
}
export default menu
