import type { AppRouteModule } from '/@/router/types'
import { LAYOUT } from '/@/router/constant'
import { t } from '/@/hooks/web/useI18n'

const genfile: AppRouteModule = {
  path: '/genfile',
  name: 'GenFile',
  component: LAYOUT,
  redirect: '/genfile/index',
  meta: {
    hideChildrenInMenu: true,
    icon: 'ion:balloon-outline',
    title: t('routes.genfile.name'),
    orderNo: 40,
  },
  children: [
    {
      path: 'index',
      name: 'GenFilePage',
      component: () => import('/@/views/page/genfile/index.vue'),
      meta: {
        title: t('routes.genfile.name'),
        icon: 'ion:balloon-outline',
        hideMenu: true,
      },
    },
  ],
}

export default genfile
