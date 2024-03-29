import type { AppRouteModule } from '/@/router/types'
import { LAYOUT } from '/@/router/constant'
import { t } from '/@/hooks/web/useI18n'

const filelist: AppRouteModule = {
  path: '/filelist',
  name: 'FileList',
  component: LAYOUT,
  redirect: '/filelist/index',
  meta: {
    hideChildrenInMenu: true,
    icon: 'ion:file-tray-full-outline',
    title: t('routes.filelist.name'),
    orderNo: 20,
  },
  children: [
    {
      path: 'index',
      name: 'FileListPage',
      component: () => import('/@/views/page/filelist/index.vue'),
      meta: {
        title: t('routes.filelist.name'),
        icon: 'ion:file-tray-full-outline',
        hideMenu: true,
      },
    },
    {
      path: 'file/:id',
      name: 'FilePage',
      component: () => import('/@/views/page/file/index.vue'),
      meta: {
        title: t('routes.filelist.file'),
        carryParam: true,
        icon: 'bi:filetype-docx',
        hideMenu: true,
      },
    },
  ],
}

export default filelist
