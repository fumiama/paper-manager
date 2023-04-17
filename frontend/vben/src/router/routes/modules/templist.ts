import type { AppRouteModule } from '/@/router/types'
import { LAYOUT } from '/@/router/constant'
import { t } from '/@/hooks/web/useI18n'

const templist: AppRouteModule = {
  path: '/templist',
  name: 'TempList',
  component: LAYOUT,
  redirect: '/templist/index',
  meta: {
    hideChildrenInMenu: true,
    icon: 'ion:ios-analytics',
    title: t('routes.templist.name'),
    orderNo: 20,
  },
  children: [
    {
      path: 'index',
      name: 'TempListPage',
      component: () => import('/@/views/page/templist/index.vue'),
      meta: {
        title: t('routes.templist.name'),
        icon: 'ion:file-tray-full-outline',
        hideMenu: true,
      },
    },
    {
      path: 'file/:id',
      name: 'TempFilePage',
      component: () => import('/@/views/page/file/index.vue'),
      meta: {
        title: t('routes.templist.file'),
        carryParam: true,
        icon: 'bi:filetype-docx',
        hideMenu: true,
      },
    },
  ],
}

export default templist
