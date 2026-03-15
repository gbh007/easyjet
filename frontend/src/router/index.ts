/**
 * router/index.ts
 *
 * Manual routes for ./src/pages/*.vue
 */

// Composables
import { createRouter, createWebHistory } from 'vue-router'
import Index from '@/pages/index.vue'
import Project from '@/pages/Project.vue'
import ProjectEditor from '@/pages/ProjectEditor.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
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
  ],
})

export default router
