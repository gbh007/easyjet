/**
 * router/index.ts
 *
 * Manual routes for ./src/pages/*.vue
 */

// Composables
import { createRouter, createWebHashHistory } from 'vue-router';
import GlobalEnvVars from '@/pages/GlobalEnvVars.vue';
import Index from '@/pages/index.vue';
import Project from '@/pages/Project.vue';
import ProjectEditor from '@/pages/ProjectEditor.vue';
import ProjectRun from '@/pages/ProjectRun.vue';

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: Index,
    },
    {
      path: '/projects/:id',
      component: Project,
    },
    {
      path: '/projects/:id/edit',
      component: ProjectEditor,
    },
    {
      path: '/projects/new',
      component: ProjectEditor,
    },
    {
      path: '/projects/:project_id/runs/:run_id',
      component: ProjectRun,
    },
    {
      path: '/global-env-vars',
      component: GlobalEnvVars,
    },
  ],
});

export default router;
